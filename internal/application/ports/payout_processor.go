package ports

import (
	"context"

	"github.com/example/payment-gateway/internal/domain"
)

// PayoutRequest captures the information required to hand off a payout to an external processor.
type PayoutRequest struct {
	WalletID     string
	ExternalID   string
	Destination  domain.PayoutDestination
	SourceAmount domain.Money
	TargetAmount domain.Money
	ExchangeRate string
	Metadata     map[string]string
}

// PayoutProcessor pushes payouts to an infrastructure layer.
type PayoutProcessor interface {
	Process(ctx context.Context, request PayoutRequest) (string, error)
}
