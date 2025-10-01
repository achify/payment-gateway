package payouts

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/example/payment-gateway/internal/application/ports"
)

// InMemoryPayoutProcessor simulates an external payout processor.
type InMemoryPayoutProcessor struct {
	counter atomic.Int64
}

// NewInMemoryPayoutProcessor creates the processor.
func NewInMemoryPayoutProcessor() *InMemoryPayoutProcessor {
	return &InMemoryPayoutProcessor{}
}

// Process records the payout and returns a synthetic identifier.
func (p *InMemoryPayoutProcessor) Process(_ context.Context, request ports.PayoutRequest) (string, error) {
	id := p.counter.Add(1)
	return fmt.Sprintf("payout-%06d", id), nil
}
