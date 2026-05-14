package wallet

import (
	"time"

	"github.com/shopspring/decimal"
)

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

type ExchangeLogItem struct {
	Id           int64           `json:"id"`
	WalletNumber string          `json:"walletNumber"`
	FromCode     string          `json:"fromCode"`
	ToCode       string          `json:"toCode"`
	FromAmount   decimal.Decimal `json:"fromAmount"`
	ToAmount     decimal.Decimal `json:"toAmount"`
	FromRate     decimal.Decimal `json:"fromRate"`  // FromCode -> KRW rate at exchange time
	FromUnit     decimal.Decimal `json:"fromUnit"`  // Unit basis for FromRate (e.g. 100 JPY)
	ToRate       decimal.Decimal `json:"toRate"`    // ToCode -> KRW rate at exchange time
	ToUnit       decimal.Decimal `json:"toUnit"`    // Unit basis for ToRate (e.g. 100 JPY)
	KrwAmount    decimal.Decimal `json:"krwAmount"` // Intermediate KRW amount
	ExchangedAt  time.Time       `json:"exchangedAt"`
}

type ExchangeLogsResponse struct {
	ExchangeLogs []ExchangeLogItem `json:"exchangeLogs"`
}
