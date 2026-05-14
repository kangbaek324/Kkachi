package wallet

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	db "github.com/kangbaek324/kkachi/db/sqlc"
	"github.com/shopspring/decimal"
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

	one := decimal.NewFromInt(1)
	fromRate, fromUnit := one, one
	toRate, toUnit := one, one

	// From -> KRW
	krwAmount := req.Amount
	if req.FromCode != "KRW" {
		rateInfo, err := q.GetRate(ctx, req.FromCode)
		if err != nil {
			return ExchangeResponse{}, fmt.Errorf("exchange: getRate %w", err)
		}

		fromRate, fromUnit = rateInfo.Rate, rateInfo.Unit
		krwAmount = req.Amount.Mul(fromRate).DivRound(fromUnit, ratePrecision)
	}

	// KRW -> To
	resultAmount := krwAmount
	if req.ToCode != "KRW" {
		rateInfo, err := q.GetRate(ctx, req.ToCode)
		if err != nil {
			return ExchangeResponse{}, fmt.Errorf("exchange: getRate %w", err)
		}
		toRate, toUnit = rateInfo.Rate, rateInfo.Unit
		resultAmount = krwAmount.Mul(toUnit).DivRound(toRate, ratePrecision)
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

	if err := q.CreateExchangeLog(ctx, db.CreateExchangeLogParams{
		WalletID:       wallet.ID,
		FromCurrencyID: from.CurrencyID,
		ToCurrencyID:   to.CurrencyID,
		FromAmount:     req.Amount,
		ToAmount:       resultAmount,
		FromRate:       fromRate,
		FromUnit:       fromUnit,
		ToRate:         toRate,
		ToUnit:         toUnit,
		KrwAmount:      krwAmount,
	}); err != nil {
		return ExchangeResponse{}, fmt.Errorf("exchange: createExchangeLog: %w", err)
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
