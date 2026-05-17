package currency

import (
	"context"
	"fmt"

	db "github.com/kangbaek324/kkachi/apps/backend/db/sqlc"
)

type Service interface {
	getCurrencies(ctx context.Context) (GetCurrenciesResponse, error)
}

type currencyService struct {
	q db.Querier
}

func NewService(q db.Querier) Service {
	return &currencyService{q: q}
}

func (s *currencyService) getCurrencies(ctx context.Context) (GetCurrenciesResponse, error) {
	currenciesDB, err := s.q.ListCurrenciesWithRate(ctx)
	if err != nil {
		return GetCurrenciesResponse{}, fmt.Errorf("getCurrencies: ListCurrencies %w", err)
	}

	var currencies []CurrencyItem
	for _, currency := range currenciesDB {
		currencies = append(currencies, CurrencyItem{
			Code: currency.Code,
			Name: currency.Name,
			Unit: currency.Unit,
			Rate: currency.Rate,
		})
	}

	return GetCurrenciesResponse{
		Currencies: currencies,
	}, nil
}
