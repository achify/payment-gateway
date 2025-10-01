package ports

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/example/payment-gateway/internal/domain"
)

// ExchangeRateProvider resolves exchange rates for supported currency pairs.
type ExchangeRateProvider interface {
	GetRate(ctx context.Context, pair domain.CurrencyPair) (decimal.Decimal, error)
}
