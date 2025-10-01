package application_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/example/payment-gateway/internal/application"
	"github.com/example/payment-gateway/internal/infrastructure/exchange"
	"github.com/example/payment-gateway/internal/infrastructure/fees"
)

func TestTransferCalculationService_Calculate(t *testing.T) {
	svc := application.NewTransferCalculationService(
		exchange.NewInMemoryExchangeRateProvider(map[string]decimal.Decimal{
			"GBP_NGN": decimal.NewFromInt(2100),
		}),
		fees.NewTieredFlatFeePolicy([]fees.FlatFeeTier{{
			Threshold: decimal.NewFromInt(50),
			Fee:       decimal.RequireFromString("1.99"),
		}}),
	)

	quote, err := svc.Calculate(context.Background(), application.TransferCalculationInput{
		SourceCurrency: "GBP",
		TargetCurrency: "NGN",
		Amount:         decimal.NewFromInt(50),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !quote.Allowed {
		t.Fatalf("expected quote to be allowed")
	}

	expectedReceive := decimal.NewFromInt(50).Mul(decimal.NewFromInt(2100))
	if !quote.ReceiveAmount.Amount.Equal(expectedReceive) {
		t.Errorf("expected receive amount %s, got %s", expectedReceive, quote.ReceiveAmount.Amount)
	}

	if quote.Fee.Amount.String() != "1.99" {
		t.Errorf("expected fee 1.99, got %s", quote.Fee.Amount.String())
	}

	expectedDebit := decimal.NewFromInt(50).Add(decimal.RequireFromString("1.99"))
	if !quote.TotalDebit.Amount.Equal(expectedDebit) {
		t.Errorf("expected total debit %s, got %s", expectedDebit, quote.TotalDebit.Amount)
	}
}

func TestTransferCalculationService_UnsupportedPair(t *testing.T) {
	svc := application.NewTransferCalculationService(
		exchange.NewInMemoryExchangeRateProvider(map[string]decimal.Decimal{}),
		fees.NewTieredFlatFeePolicy(nil),
	)

	_, err := svc.Calculate(context.Background(), application.TransferCalculationInput{
		SourceCurrency: "GBP",
		TargetCurrency: "GBP",
		Amount:         decimal.NewFromInt(10),
	})
	if err == nil {
		t.Fatalf("expected error for unsupported pair")
	}
}
