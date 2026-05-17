package wallet

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	db "github.com/kangbaek324/kkachi/apps/backend/db/sqlc"
)

// Lock the wallet row first to serialize concurrent transfers on the same wallet.
func (s *walletService) Transfer(ctx context.Context, req TransferRequest, walletNumber string, userId int64) error {
	if walletNumber == req.Receiver {
		return ErrSelfTransfer
	}
	if req.Amount.LessThan(minAmount) {
		return ErrInvalidAmount
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
	if _, err := q.UpsertBalance(ctx, db.UpsertBalanceParams{
		WalletID:   sender.WalletID,
		CurrencyID: sender.CurrencyID,
		Amount:     req.Amount.Neg(),
	}); err != nil {
		return fmt.Errorf("transfer: upsertBalance: %w", err)
	}
	if _, err := q.UpsertBalance(ctx, db.UpsertBalanceParams{
		WalletID:   receiver.WalletID,
		CurrencyID: receiver.CurrencyID,
		Amount:     req.Amount,
	}); err != nil {
		return fmt.Errorf("transfer: upsertBalance:: %w", err)
	}

	if err := q.CreateTransferLog(ctx, db.CreateTransferLogParams{
		FromWalletID: sender.WalletID,
		ToWalletID:   receiver.WalletID,
		CurrencyID:   sender.CurrencyID,
		Amount:       req.Amount,
	}); err != nil {
		return fmt.Errorf("transfer: createTransferLog: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transfer: commit: %w", err)
	}
	return nil
}

func (s *walletService) GetTransferLogs(ctx context.Context, walletNumber string, userId int64) (TransferLogsResponse, error) {
	wallet, err := s.q.GetWallet(ctx, walletNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TransferLogsResponse{}, ErrWalletNotFound
		}
		return TransferLogsResponse{}, fmt.Errorf("getTransferLogs: getWallet: %w", err)
	}
	if wallet.UserID != userId {
		return TransferLogsResponse{}, ErrNotWalletOwner
	}

	rows, err := s.q.GetTransferLogs(ctx, wallet.ID)
	if err != nil {
		return TransferLogsResponse{}, fmt.Errorf("getTransferLogs: %w", err)
	}

	items := make([]TransferLogItem, len(rows))
	for i, r := range rows {
		items[i] = TransferLogItem{
			Id:               r.ID,
			FromWalletNumber: r.FromWalletNumber,
			ToWalletNumber:   r.ToWalletNumber,
			Amount:           r.Amount,
			CurrencyCode:     r.CurrencyCode,
			TransferredAt:    r.TransferredAt.Time,
		}
	}

	return TransferLogsResponse{TransferLogs: items}, nil
}
