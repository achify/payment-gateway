package exchange

import (
	"context"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/example/payment-gateway/internal/domain"
)

// InMemoryExchangeRateProvider provides deterministic exchange rates for testing and demos.
type InMemoryExchangeRateProvider struct {
	rates map[string]decimal.Decimal
}

// NewInMemoryExchangeRateProvider creates a provider with the supplied rates.
func NewInMemoryExchangeRateProvider(rates map[string]decimal.Decimal) *InMemoryExchangeRateProvider {
	normalised := make(map[string]decimal.Decimal, len(rates))
	for k, v := range rates {
		normalised[strings.ToUpper(k)] = v
	}
	return &InMemoryExchangeRateProvider{rates: normalised}
}

func pairKey(pair domain.CurrencyPair) string {
	return fmt.Sprintf("%s_%s", pair.Base, pair.Quote)
}

// GetRate returns the configured rate or an error when unsupported.
func (p *InMemoryExchangeRateProvider) GetRate(_ context.Context, pair domain.CurrencyPair) (decimal.Decimal, error) {
	rate, ok := p.rates[pairKey(pair)]
	if !ok {
		return decimal.Zero, fmt.Errorf("exchange rate for pair %s/%s not found", pair.Base, pair.Quote)
	}
	return rate, nil
}
