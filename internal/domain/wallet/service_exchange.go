package wallet

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	db "github.com/kangbaek324/kkachi/db/sqlc"
)

func (s *walletService) Exchange(ctx context.Context, req ExchangeRequest, walletNumber string, userId int64) (ExchangeResponse, error) {
	if req.FromCode == req.ToCode {
		return ExchangeResponse{}, ErrSameCurrency
	}
	if req.Amount.LessThan(minAmount) {
		return ExchangeResponse{}, ErrInvalidAmount
	}

	wallet, err := s.q.GetWallet(ctx, walletNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ExchangeResponse{}, ErrWalletNotFound
		}
		return ExchangeResponse{}, fmt.Errorf("exchange: %w", err)
	}

	if wallet.UserID != userId {
		return ExchangeResponse{}, ErrNotWalletOwner
	}

	// Processing Exchange
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return ExchangeResponse{}, fmt.Errorf("exchange: openTx: %w", err)
	}
	defer tx.Rollback(ctx)

	q := db.New(tx)

	// Find FromCode Balance
	from, err := q.GetWalletBalanceLock(ctx, db.GetWalletBalanceLockParams{
		WalletNumber: walletNumber,
		Code:         req.FromCode,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ExchangeResponse{}, ErrInsufficientFund
		}
		return ExchangeResponse{}, fmt.Errorf("exchange: getWalletBalanceLock %w", err)
	}
	if from.Amount.LessThan(req.Amount) {
		return ExchangeResponse{}, ErrInsufficientFund
	}

	// Find toCode Balance
	to, err := q.GetWalletBalance(ctx, db.GetWalletBalanceParams{
		WalletNumber: walletNumber,
		Code:         req.ToCode,
	})
	if err != nil {
		return ExchangeResponse{}, fmt.Errorf("exchange: getWalletBalanceLock %w", err)
	}

	const ratePrecision = 8
	const amountPrecision = 6

	krwAmount := req.Amount
	if req.FromCode != "KRW" {
		rateInfo, err := q.GetRate(ctx, req.FromCode)
		if err != nil {
			return ExchangeResponse{}, fmt.Errorf("exchange: getRate %w", err)
		}

		krwAmount = req.Amount.Mul(rateInfo.Rate).DivRound(rateInfo.Unit, ratePrecision)
	}

	resultAmount := krwAmount
	if req.ToCode != "KRW" {
		rateInfo, err := q.GetRate(ctx, req.ToCode)
		if err != nil {
			return ExchangeResponse{}, fmt.Errorf("exchange: getRate %w", err)
		}
		resultAmount = krwAmount.Mul(rateInfo.Unit).DivRound(rateInfo.Rate, ratePrecision)
	}
	resultAmount = resultAmount.Round(amountPrecision)

	// Upsert Result
	fromBalance, err := q.UpsertBalance(ctx, db.UpsertBalanceParams{
		WalletID:   wallet.ID,
		CurrencyID: from.CurrencyID,
		Amount:     req.Amount.Neg(),
	})
	if err != nil {
		return ExchangeResponse{}, fmt.Errorf("exchange: upsertBalance: %w", err)
	}

	toBalance, err := q.UpsertBalance(ctx, db.UpsertBalanceParams{
		WalletID:   wallet.ID,
		CurrencyID: to.CurrencyID,
		Amount:     resultAmount,
	})
	if err != nil {
		return ExchangeResponse{}, fmt.Errorf("exchange: upsertBalance: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return ExchangeResponse{}, fmt.Errorf("exchange: commit: %w", err)
	}

	return ExchangeResponse{
		FromCode:        req.FromCode,
		ToCode:          req.ToCode,
		ConvertedAmount: resultAmount,
		FromBalance:     fromBalance,
		ToBalance:       toBalance,
	}, nil
}
