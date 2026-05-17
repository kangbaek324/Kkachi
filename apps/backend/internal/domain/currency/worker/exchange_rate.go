package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	db "github.com/kangbaek324/kkachi/apps/backend/db/sqlc"
	"github.com/robfig/cron/v3"
	"github.com/shopspring/decimal"
)

type exchnageRateWorker struct {
	q db.Querier
}

func NewExchangeRateWorker(q db.Querier) *exchnageRateWorker {
	return &exchnageRateWorker{q: q}
}

type apiResponse struct {
	CurUnit  string `json:"cur_unit"`
	DealBasR string `json:"deal_bas_r"`
}

func (w *exchnageRateWorker) Start(apiKey string) {
	c := cron.New(cron.WithLocation(time.FixedZone("KST", 9*60*60)))
	// every day at 12:00 PM (KST)
	c.AddFunc("0 12 * * *", func() {
		currenciesRow, err := w.q.ListCurrencies(context.Background())
		if err != nil {
			log.Printf("exchangeRateWorker: getCurrenciesRow: %v", err)
			return
		}

		// Currency Code Setup
		currencyMap := map[string]int64{}
		for _, cr := range currenciesRow {
			currencyMap[cr.Code] = cr.ID
		}

		// Rates Setting
		rates, err := fetchRates(apiKey)
		if err != nil {
			log.Printf("exchangeRateWorker: fetchRates: %v", err)
			return
		}

		for _, apiRate := range rates {
			rateStr := strings.ReplaceAll(apiRate.DealBasR, ",", "")
			rateDec, err := decimal.NewFromString(rateStr)
			if err != nil {
				log.Printf("exchangeRateWorker: parse rate %s: %v", apiRate.CurUnit, err)
				continue
			}

			// DB Upsert
			if err := w.q.UpsertExchangeRate(context.Background(), db.UpsertExchangeRateParams{
				CurrencyID: currencyMap[apiRate.CurUnit],
				Rate:       rateDec,
			}); err != nil {
				log.Printf("exchangeRateWorker: upsert %s: %v", apiRate.CurUnit, err)
				continue
			}
		}
	})

	c.Start()
}

func fetchRates(apiKey string) ([]apiResponse, error) {
	url := fmt.Sprintf(
		"https://oapi.koreaexim.go.kr/site/program/financial/exchangeJSON?authkey=%s&searchdate=%s&data=AP01",
		apiKey,
		"20260512",
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rates []apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return nil, err
	}

	return rates, nil
}
