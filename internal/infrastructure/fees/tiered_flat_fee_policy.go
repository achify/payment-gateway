package fees

import (
	"context"
	"sort"

	"github.com/shopspring/decimal"

	"github.com/example/payment-gateway/internal/domain"
)

// FlatFeeTier represents a fee rule applied up to and including the threshold amount.
type FlatFeeTier struct {
	Threshold decimal.Decimal
	Fee       decimal.Decimal
}

// TieredFlatFeePolicy applies flat fees based on configured tiers.
type TieredFlatFeePolicy struct {
	tiers []FlatFeeTier
}

// NewTieredFlatFeePolicy builds a policy sorted by threshold ascending.
func NewTieredFlatFeePolicy(tiers []FlatFeeTier) *TieredFlatFeePolicy {
	sorted := make([]FlatFeeTier, len(tiers))
	copy(sorted, tiers)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Threshold.LessThan(sorted[j].Threshold)
	})
	return &TieredFlatFeePolicy{tiers: sorted}
}

// CalculateFee returns the matching flat fee or zero when no tier applies.
func (p *TieredFlatFeePolicy) CalculateFee(_ context.Context, amount domain.Money) (domain.Money, error) {
	for _, tier := range p.tiers {
		if amount.Amount.LessThanOrEqual(tier.Threshold) {
			feeMoney, err := domain.NewMoney(amount.Currency, tier.Fee)
			if err != nil {
				return domain.Money{}, err
			}
			return feeMoney, nil
		}
	}
	return domain.ZeroMoney(amount.Currency), nil
}
