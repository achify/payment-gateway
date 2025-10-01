package main

import (
	"log"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/example/payment-gateway/internal/application"
	"github.com/example/payment-gateway/internal/domain"
	internalhttp "github.com/example/payment-gateway/internal/http"
	infrEvents "github.com/example/payment-gateway/internal/infrastructure/events"
	infrExchange "github.com/example/payment-gateway/internal/infrastructure/exchange"
	infrFees "github.com/example/payment-gateway/internal/infrastructure/fees"
	infrPayouts "github.com/example/payment-gateway/internal/infrastructure/payouts"
	infrRepos "github.com/example/payment-gateway/internal/infrastructure/repositories"
)

func main() {
	rateProvider := infrExchange.NewInMemoryExchangeRateProvider(map[string]decimal.Decimal{
		"GBP_NGN": decimal.NewFromInt(2100),
	})

	feePolicy := infrFees.NewTieredFlatFeePolicy([]infrFees.FlatFeeTier{
		{Threshold: decimal.NewFromInt(50), Fee: decimal.RequireFromString("1.99")},
	})

	balance, err := domain.NewMoney("GBP", decimal.NewFromInt(100))
	if err != nil {
		log.Fatalf("failed to create initial balance: %v", err)
	}

	walletRepo := infrRepos.NewInMemoryWalletRepository(map[string]*domain.Wallet{
		"wallet-123": {
			ID:       "wallet-123",
			Currency: "GBP",
			Balance:  balance,
		},
	})

	payoutProcessor := infrPayouts.NewInMemoryPayoutProcessor()
	eventPublisher := infrEvents.NewLoggingEventPublisher()

	transferService := application.NewTransferCalculationService(rateProvider, feePolicy)
	payoutService := application.NewPayoutService(walletRepo, payoutProcessor, eventPublisher, transferService)

	server := internalhttp.NewServer(transferService, payoutService)

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", server.Router()); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
