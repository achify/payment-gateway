package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/example/payment-gateway/internal/application/ports"
	"github.com/example/payment-gateway/internal/domain"
)

// PayoutInput contains the data required to initiate a payout.
type PayoutInput struct {
	WalletID       string
	Amount         decimal.Decimal
	SourceCurrency string
	TargetCurrency string
	Destination    domain.PayoutDestination
	Reference      string
}

// PayoutResult summarises the successful payout initiation.
type PayoutResult struct {
	PayoutID         string
	Quote            domain.TransferQuote
	RemainingBalance domain.Money
}

// PayoutService handles payout orchestration while enforcing domain rules.
type PayoutService struct {
	wallets         ports.WalletRepository
	payoutProcessor ports.PayoutProcessor
	eventPublisher  ports.EventPublisher
	transferService *TransferCalculationService
}

// NewPayoutService constructs a PayoutService.
func NewPayoutService(wallets ports.WalletRepository, payoutProcessor ports.PayoutProcessor, eventPublisher ports.EventPublisher, transferService *TransferCalculationService) *PayoutService {
	return &PayoutService{
		wallets:         wallets,
		payoutProcessor: payoutProcessor,
		eventPublisher:  eventPublisher,
		transferService: transferService,
	}
}

// InitiatePayout validates the request, debits the wallet, and triggers the payout.
func (s *PayoutService) InitiatePayout(ctx context.Context, input PayoutInput) (PayoutResult, error) {
	if input.WalletID == "" {
		return PayoutResult{}, errors.New("wallet ID is required")
	}
	if input.Reference == "" {
		input.Reference = uuid.NewString()
	}

	wallet, err := s.wallets.GetByID(ctx, input.WalletID)
	if err != nil {
		return PayoutResult{}, err
	}

	if wallet.Currency != input.SourceCurrency {
		return PayoutResult{}, domain.ErrCurrencyMismatchWallet
	}

	quote, err := s.transferService.Calculate(ctx, TransferCalculationInput{
		SourceCurrency: input.SourceCurrency,
		TargetCurrency: input.TargetCurrency,
		Amount:         input.Amount,
	})
	if err != nil {
		return PayoutResult{}, err
	}

	// Ensure the wallet has sufficient funds including fees.
	if err := wallet.Debit(quote.TotalDebit); err != nil {
		return PayoutResult{}, err
	}

	if err := s.wallets.Save(ctx, wallet); err != nil {
		return PayoutResult{}, err
	}

	request := ports.PayoutRequest{
		WalletID:     wallet.ID,
		ExternalID:   input.Reference,
		Destination:  input.Destination,
		SourceAmount: quote.SendAmount,
		TargetAmount: quote.ReceiveAmount,
		ExchangeRate: quote.ExchangeRate.String(),
		Metadata: map[string]string{
			"fee":        quote.Fee.Amount.String(),
			"totalDebit": quote.TotalDebit.Amount.String(),
		},
	}

	payoutID, err := s.payoutProcessor.Process(ctx, request)
	if err != nil {
		return PayoutResult{}, err
	}

	_ = s.eventPublisher.Publish(ctx, "payout.initiated", map[string]any{
		"payoutId":       payoutID,
		"walletId":       wallet.ID,
		"reference":      input.Reference,
		"sourceCurrency": quote.Pair.Base,
		"targetCurrency": quote.Pair.Quote,
		"amount":         quote.SendAmount.Amount.String(),
	})

	return PayoutResult{
		PayoutID:         payoutID,
		Quote:            quote,
		RemainingBalance: wallet.Balance,
	}, nil
}
