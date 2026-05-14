package currency

import "github.com/shopspring/decimal"

type CurrencyItem struct {
	Code string              `json:"code"`
	Name string              `json:"name"`
	Unit decimal.Decimal     `json:"unit"`
	Rate decimal.NullDecimal `json:"rate"`
}

type GetCurrenciesResponse struct {
	Currencies []CurrencyItem `json:"currencies"`
}
