package wallet

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransferRequest struct {
	Receiver     string          `json:"receiver" binding:"required"`
	CurrencyCode string          `json:"currencyCode" binding:"required"`
	Amount       decimal.Decimal `json:"amount" binding:"required"`
}

type TransferLogItem struct {
	FromWalletNumber string          `json:"fromWalletNumber"`
	ToWalletNumber   string          `json:"toWalletNumber"`
	Amount           decimal.Decimal `json:"amount"`
	CurrencyCode     string          `json:"currencyCode"`
	TransferredAt    time.Time       `json:"transferredAt"`
}

type TransferLogsResponse struct {
	TransferLogs []TransferLogItem `json:"transferLogs"`
}
