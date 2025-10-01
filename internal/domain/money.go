package domain

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidAmount    = errors.New("amount must be greater than zero")
	ErrCurrencyMismatch = errors.New("currency mismatch")
)

// Money represents a monetary value in a specific currency using arbitrary precision decimals.
type Money struct {
	Currency string
	Amount   decimal.Decimal
}

// NewMoney creates a new Money instance ensuring the amount is positive.
func NewMoney(currency string, amount decimal.Decimal) (Money, error) {
	if currency == "" {
		return Money{}, errors.New("currency is required")
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return Money{}, ErrInvalidAmount
	}
	return Money{Currency: currency, Amount: amount}, nil
}

// Add returns a new Money that is the sum of m and other.
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, ErrCurrencyMismatch
	}
	return Money{Currency: m.Currency, Amount: m.Amount.Add(other.Amount)}, nil
}

// Subtract returns a new Money that is the result of subtracting other from m.
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, ErrCurrencyMismatch
	}
	result := m.Amount.Sub(other.Amount)
	if result.IsNegative() {
		return Money{}, ErrInvalidAmount
	}
	return Money{Currency: m.Currency, Amount: result}, nil
}

// Multiply returns a new Money that is the result of multiplying m.Amount by factor.
func (m Money) Multiply(factor decimal.Decimal) Money {
	return Money{Currency: m.Currency, Amount: m.Amount.Mul(factor)}
}

// ZeroMoney returns a zero value for the provided currency.
func ZeroMoney(currency string) Money {
	return Money{Currency: currency, Amount: decimal.Zero}
}
