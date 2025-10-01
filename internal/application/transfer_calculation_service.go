package application

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"

	"github.com/example/payment-gateway/internal/application/ports"
	"github.com/example/payment-gateway/internal/domain"
)

// TransferCalculationInput defines the payload for computing a transfer quote.
type TransferCalculationInput struct {
	SourceCurrency string
	TargetCurrency string
	Amount         decimal.Decimal
}

// TransferCalculationService provides transfer quotes using the configured policies.
type TransferCalculationService struct {
	exchangeRateProvider ports.ExchangeRateProvider
	feePolicy            ports.FeePolicy
}

// NewTransferCalculationService creates a new service instance.
func NewTransferCalculationService(exchangeRateProvider ports.ExchangeRateProvider, feePolicy ports.FeePolicy) *TransferCalculationService {
	return &TransferCalculationService{
		exchangeRateProvider: exchangeRateProvider,
		feePolicy:            feePolicy,
	}
}

// Calculate produces a transfer quote from the supplied input.
func (s *TransferCalculationService) Calculate(ctx context.Context, input TransferCalculationInput) (domain.TransferQuote, error) {
	pair := domain.NewCurrencyPair(input.SourceCurrency, input.TargetCurrency)
	if !pair.IsValid() {
		return domain.TransferQuote{Allowed: false, Pair: pair}, errors.New("unsupported currency pair")
	}

	sendMoney, err := domain.NewMoney(pair.Base, input.Amount)
	if err != nil {
		return domain.TransferQuote{Allowed: false, Pair: pair}, err
	}

	rate, err := s.exchangeRateProvider.GetRate(ctx, pair)
	if err != nil {
		return domain.TransferQuote{Allowed: false, Pair: pair}, err
	}

	receiveAmount := domain.Money{Currency: pair.Quote, Amount: sendMoney.Amount.Mul(rate)}

	fee, err := s.feePolicy.CalculateFee(ctx, sendMoney)
	if err != nil {
		return domain.TransferQuote{Allowed: false, Pair: pair}, err
	}

	totalDebit, err := sendMoney.Add(fee)
	if err != nil {
		return domain.TransferQuote{Allowed: false, Pair: pair}, err
	}

	return domain.TransferQuote{
		Pair:          pair,
		ExchangeRate:  rate,
		SendAmount:    sendMoney,
		ReceiveAmount: receiveAmount,
		Fee:           fee,
		TotalDebit:    totalDebit,
		Allowed:       true,
	}, nil
}
