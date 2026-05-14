package wallet

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	db "github.com/kangbaek324/kkachi/db/sqlc"
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

	return tx.Commit(ctx)
}
