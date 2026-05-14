package wallet

import (
	"github.com/shopspring/decimal"
)

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
