package ports

import (
	"context"

	"github.com/example/payment-gateway/internal/domain"
)

// FeePolicy determines the fee to apply for a given transfer.
type FeePolicy interface {
	CalculateFee(ctx context.Context, amount domain.Money) (domain.Money, error)
}
