package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewMoney(t *testing.T) {
	amount := decimal.NewFromFloat(100.50)
	money, err := NewMoney("USD", amount)
	assert.NoError(t, err)
	assert.Equal(t, "USD", money.Currency)
	assert.True(t, money.Amount.Equal(amount))
}

func TestNewMoneyValidation(t *testing.T) {
	_, err := NewMoney("", decimal.NewFromInt(10))
	assert.Error(t, err)
	_, err = NewMoney("USD", decimal.Zero)
	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestMoneyAdd(t *testing.T) {
	base, _ := NewMoney("USD", decimal.NewFromInt(50))
	other, _ := NewMoney("USD", decimal.NewFromInt(25))

	result, err := base.Add(other)
	assert.NoError(t, err)
	expected := decimal.NewFromInt(75)
	assert.True(t, result.Amount.Equal(expected))
}

func TestMoneyAddCurrencyMismatch(t *testing.T) {
	base, _ := NewMoney("USD", decimal.NewFromInt(50))
	other, _ := NewMoney("EUR", decimal.NewFromInt(25))

	_, err := base.Add(other)
	assert.ErrorIs(t, err, ErrCurrencyMismatch)
}

func TestMoneySubtract(t *testing.T) {
	base, _ := NewMoney("USD", decimal.NewFromInt(100))
	other, _ := NewMoney("USD", decimal.NewFromInt(30))

	result, err := base.Subtract(other)
	assert.NoError(t, err)
	expected := decimal.NewFromInt(70)
	assert.True(t, result.Amount.Equal(expected))
}

func TestMoneySubtractValidation(t *testing.T) {
	base, _ := NewMoney("USD", decimal.NewFromInt(30))
	larger, _ := NewMoney("USD", decimal.NewFromInt(40))
	otherCurrency, _ := NewMoney("EUR", decimal.NewFromInt(10))

	_, err := base.Subtract(larger)
	assert.ErrorIs(t, err, ErrInvalidAmount)
	_, err = base.Subtract(otherCurrency)
	assert.ErrorIs(t, err, ErrCurrencyMismatch)
}

func TestMoneyMultiply(t *testing.T) {
	base, _ := NewMoney("USD", decimal.NewFromFloat(10.25))
	factor := decimal.NewFromFloat(1.5)

	result := base.Multiply(factor)
	expected := decimal.RequireFromString("15.375")
	assert.True(t, result.Amount.Equal(expected))
	assert.Equal(t, base.Currency, result.Currency)
}

func TestZeroMoney(t *testing.T) {
	zero := ZeroMoney("GBP")
	assert.Equal(t, "GBP", zero.Currency)
	assert.True(t, zero.Amount.Equal(decimal.Zero))
}
