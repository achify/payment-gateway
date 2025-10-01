package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/example/payment-gateway/internal/application"
	"github.com/example/payment-gateway/internal/application/ports"
	"github.com/example/payment-gateway/internal/domain"
	"github.com/example/payment-gateway/internal/infrastructure/exchange"
	"github.com/example/payment-gateway/internal/infrastructure/fees"
	"github.com/example/payment-gateway/internal/infrastructure/payouts"
	"github.com/example/payment-gateway/internal/infrastructure/repositories"
)

type stubEventPublisher struct {
	published bool
}

func (s *stubEventPublisher) Publish(ctx context.Context, topic string, payload any) error {
	s.published = true
	return nil
}

func setupPayoutService(t *testing.T) (*application.PayoutService, *repositories.InMemoryWalletRepository, *stubEventPublisher) {
	t.Helper()

	transferService := application.NewTransferCalculationService(
		exchange.NewInMemoryExchangeRateProvider(map[string]decimal.Decimal{
			"GBP_NGN": decimal.NewFromInt(2100),
		}),
		fees.NewTieredFlatFeePolicy([]fees.FlatFeeTier{{
			Threshold: decimal.NewFromInt(50),
			Fee:       decimal.RequireFromString("1.99"),
		}}),
	)

	balance, _ := domain.NewMoney("GBP", decimal.NewFromInt(100))
	walletRepo := repositories.NewInMemoryWalletRepository(map[string]*domain.Wallet{
		"wallet-123": {ID: "wallet-123", Currency: "GBP", Balance: balance},
	})

	eventPublisher := &stubEventPublisher{}

	service := application.NewPayoutService(
		walletRepo,
		payouts.NewInMemoryPayoutProcessor(),
		eventPublisher,
		transferService,
	)

	return service, walletRepo, eventPublisher
}

func TestPayoutService_InitiatePayout(t *testing.T) {
	service, _, events := setupPayoutService(t)

	result, err := service.InitiatePayout(context.Background(), application.PayoutInput{
		WalletID:       "wallet-123",
		Amount:         decimal.NewFromInt(50),
		SourceCurrency: "GBP",
		TargetCurrency: "NGN",
		Destination: domain.PayoutDestination{
			Type:          "bank_account",
			AccountNumber: "1234567890",
			BankCode:      "999",
			Country:       "NG",
			Currency:      "NGN",
			RecipientName: "John Doe",
		},
		Reference: "client-ref-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.PayoutID == "" {
		t.Fatalf("expected payout ID to be generated")
	}

	expectedBalance := decimal.NewFromInt(100).Sub(decimal.NewFromInt(50).Add(decimal.RequireFromString("1.99")))
	if !result.RemainingBalance.Amount.Equal(expectedBalance) {
		t.Fatalf("expected remaining balance %s, got %s", expectedBalance, result.RemainingBalance.Amount)
	}

	if !events.published {
		t.Fatalf("expected event to be published")
	}
}

type failingPayoutProcessor struct{}

func (f *failingPayoutProcessor) Process(ctx context.Context, request ports.PayoutRequest) (string, error) {
	return "", errors.New("processor unavailable")
}

func TestPayoutService_PayoutProcessorFailure(t *testing.T) {
	transferService := application.NewTransferCalculationService(
		exchange.NewInMemoryExchangeRateProvider(map[string]decimal.Decimal{
			"GBP_NGN": decimal.NewFromInt(2100),
		}),
		fees.NewTieredFlatFeePolicy([]fees.FlatFeeTier{{
			Threshold: decimal.NewFromInt(50),
			Fee:       decimal.RequireFromString("1.99"),
		}}),
	)

	balance, _ := domain.NewMoney("GBP", decimal.NewFromInt(100))
	walletRepo := repositories.NewInMemoryWalletRepository(map[string]*domain.Wallet{
		"wallet-123": {ID: "wallet-123", Currency: "GBP", Balance: balance},
	})

	eventPublisher := &stubEventPublisher{}

	service := application.NewPayoutService(
		walletRepo,
		&failingPayoutProcessor{},
		eventPublisher,
		transferService,
	)

	_, err := service.InitiatePayout(context.Background(), application.PayoutInput{
		WalletID:       "wallet-123",
		Amount:         decimal.NewFromInt(10),
		SourceCurrency: "GBP",
		TargetCurrency: "NGN",
		Destination:    domain.PayoutDestination{Currency: "NGN"},
	})
	if err == nil {
		t.Fatalf("expected processor failure error")
	}
}

func TestPayoutService_InsufficientFunds(t *testing.T) {
	service, _, _ := setupPayoutService(t)

	_, err := service.InitiatePayout(context.Background(), application.PayoutInput{
		WalletID:       "wallet-123",
		Amount:         decimal.NewFromInt(200),
		SourceCurrency: "GBP",
		TargetCurrency: "NGN",
		Destination:    domain.PayoutDestination{Currency: "NGN"},
	})
	if err == nil {
		t.Fatalf("expected insufficient funds error")
	}
}
