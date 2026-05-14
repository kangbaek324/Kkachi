package wallet

import "github.com/shopspring/decimal"

type CreateWalletRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

type CreateWalletResponse struct {
	WalletNumber string `json:"walletNumber"`
	Nickname     string `json:"nickname"`
}

type WalletItem struct {
	WalletNumber string `json:"walletNumber"`
	Nickname     string `json:"nickname"`
}

type GetWalletsResponse struct {
	Wallets []WalletItem `json:"wallets"`
}

type WalletBalance struct {
	Code   string          `json:"code"`
	Name   string          `json:"name"`
	Amount decimal.Decimal `json:"amount"`
}

type GetWalletBalanceResponse struct {
	Balances []WalletBalance `json:"balances"`
}

type EditWalletNicknameRequest struct {
	WalletNumber string `json:"walletNumber" binding:"required"`
	Nickname     string `json:"nickname"`
}

// Transfer
type TransferRequest struct {
	Receiver     string          `json:"receiver" binding:"required"`
	CurrencyCode string          `json:"currencyCode" binding:"required"`
	Amount       decimal.Decimal `json:"amount" binding:"required"`
}

// Exchange
type ExchangeRequest struct {
	FromCode string          `json:"fromCode" binding:"required"`
	ToCode   string          `json:"toCode" binding:"required"`
	Amount   decimal.Decimal `json:"amount" binding:"required"`
}

type ExchangeResponse struct {
	FromCode        string          `json:"fromCode"`
	ToCode          string          `json:"toCode"`
	ConvertedAmount decimal.Decimal `json:"convertedAmount"`
	FromBalance     decimal.Decimal `json:"fromBalance"`
	ToBalance       decimal.Decimal `json:"toBalance"`
}
