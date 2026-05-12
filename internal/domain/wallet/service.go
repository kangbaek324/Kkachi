package wallet

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/kangbaek324/kkachi/db/sqlc"
	"github.com/kangbaek324/kkachi/internal/common"
)

var ErrWalletNotFound = common.NewAppError(http.StatusNotFound, "wallet not found")
var ErrNotWalletOwner = common.NewAppError(http.StatusForbidden, "not wallet owner")
var ErrInsufficientFund = common.NewAppError(http.StatusBadRequest, "insufficient funds")
var ErrReceiverNotFound = common.NewAppError(http.StatusNotFound, "receiver wallet not found")
var ErrSelfTransfer = common.NewAppError(http.StatusBadRequest, "cannot transfer to the same wallet")

type Service interface {
	CreateWallet(ctx context.Context, req CreateWalletRequest, userId int64) (CreateWalletResponse, error)
	GetWallets(ctx context.Context, userId int64) (GetWalletsResponse, error)
	EditWalletNickname(ctx context.Context, req EditWalletNicknameRequest, userId int64) error
	GetWalletBalances(ctx context.Context, userId int64, walletNumber string) (GetWalletBalanceResponse, error)
	Transfer(ctx context.Context, req TransferRequest, walletNumber string, userId int64) error
}

type walletService struct {
	q    db.Querier
	pool *pgxpool.Pool
}

func NewService(q db.Querier, pool *pgxpool.Pool) Service {
	return &walletService{q: q, pool: pool}
}

func (s *walletService) CreateWallet(ctx context.Context, req CreateWalletRequest, userId int64) (CreateWalletResponse, error) {
	for range 5 {
		walletNumber := fmt.Sprintf("KK-%06d", rand.IntN(900000)+100000)
		wallet, err := s.q.CreateWallet(ctx, db.CreateWalletParams{
			UserID:       userId,
			WalletNumber: walletNumber,
			Nickname:     req.Nickname,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				continue
			}
			return CreateWalletResponse{}, fmt.Errorf("createWallet: %w", err)
		}

		return CreateWalletResponse{
			WalletNumber: wallet.WalletNumber,
			Nickname:     wallet.Nickname,
		}, nil
	}

	return CreateWalletResponse{}, fmt.Errorf("createWallet: failed to generate unique wallet number")
}

func (s *walletService) GetWallets(ctx context.Context, userId int64) (GetWalletsResponse, error) {
	wallets, err := s.q.GetWallets(ctx, userId)
	if err != nil {
		return GetWalletsResponse{}, fmt.Errorf("getWallets: %w", err)
	}

	items := make([]WalletItem, len(wallets))
	for i, w := range wallets {
		items[i] = WalletItem{
			WalletNumber: w.WalletNumber,
			Nickname:     w.Nickname,
		}
	}

	return GetWalletsResponse{
		Wallets: items,
	}, nil
}

func (s *walletService) GetWalletBalances(ctx context.Context, userId int64, walletNumber string) (GetWalletBalanceResponse, error) {
	rows, err := s.q.GetWalletBalances(ctx, walletNumber)
	if err != nil {
		return GetWalletBalanceResponse{}, fmt.Errorf("GetWalletBalance: %w", err)
	}

	if len(rows) == 0 {
		return GetWalletBalanceResponse{}, ErrWalletNotFound
	}
	if rows[0].UserID != userId {
		return GetWalletBalanceResponse{}, ErrNotWalletOwner
	}

	items := make([]WalletBalance, len(rows))
	for i, r := range rows {
		items[i] = WalletBalance{
			Code:   r.Code,
			Name:   r.Name,
			Amount: r.Amount,
		}
	}

	return GetWalletBalanceResponse{
		Balances: items,
	}, nil
}

func (s *walletService) EditWalletNickname(ctx context.Context, req EditWalletNicknameRequest, userId int64) error {
	wallet, err := s.q.GetWallet(ctx, req.WalletNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrWalletNotFound
		}
		return fmt.Errorf("editWalletNickname: %w", err)
	}

	if wallet.UserID != userId {
		return ErrNotWalletOwner
	}

	return nil
}

// Lock the wallet row first to serialize concurrent transfers on the same wallet.
func (s *walletService) Transfer(ctx context.Context, req TransferRequest, walletNumber string, userId int64) error {
	if walletNumber == req.Receiver {
		return ErrSelfTransfer
	}

	// Validate wallet
	wallet, err := s.q.GetWallet(ctx, walletNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrWalletNotFound
		}
		return fmt.Errorf("transfer: getWallet: %w", err)
	}
	if wallet.UserID != userId {
		return ErrNotWalletOwner
	}

	// Processing Transfer
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("transfer: openTx: %w", err)
	}
	defer tx.Rollback(ctx)

	q := db.New(tx)

	// Find Sender Balance
	sender, err := q.GetWalletBalanceLock(ctx, db.GetWalletBalanceLockParams{
		WalletNumber: walletNumber,
		Code:         req.CurrencyCode,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInsufficientFund
		}
		return fmt.Errorf("transfer: getWalletBalanceLock %w", err)
	}
	if sender.Amount.LessThan(req.Amount) {
		return ErrInsufficientFund
	}

	// Find Receiver Balance
	receiver, err := q.GetWalletBalanceLock(ctx, db.GetWalletBalanceLockParams{
		WalletNumber: req.Receiver,
		Code:         req.CurrencyCode,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrReceiverNotFound
		}
		return fmt.Errorf("transfer: getWalletBalanceLock: %w", err)
	}

	// Update Result
	if err := q.UpsertBalance(ctx, db.UpsertBalanceParams{
		WalletID:   sender.WalletID,
		CurrencyID: sender.CurrencyID,
		Amount:     req.Amount.Neg(),
	}); err != nil {
		return fmt.Errorf("transfer: upsertBalance: %w", err)
	}
	if err := q.UpsertBalance(ctx, db.UpsertBalanceParams{
		WalletID:   receiver.WalletID,
		CurrencyID: receiver.CurrencyID,
		Amount:     req.Amount,
	}); err != nil {
		return fmt.Errorf("transfer: upsertBalance:: %w", err)
	}

	return tx.Commit(ctx)
}
