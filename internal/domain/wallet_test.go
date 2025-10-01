package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestWalletDebit(t *testing.T) {
	balance, _ := NewMoney("USD", decimal.NewFromInt(100))
	wallet := Wallet{ID: "w1", Currency: "USD", Balance: balance}
	debit, _ := NewMoney("USD", decimal.NewFromInt(60))

	err := wallet.Debit(debit)
	assert.NoError(t, err)
	expected := decimal.NewFromInt(40)
	assert.True(t, wallet.Balance.Amount.Equal(expected))
}

func TestWalletDebitValidation(t *testing.T) {
	balance, _ := NewMoney("USD", decimal.NewFromInt(50))
	wallet := Wallet{ID: "w1", Currency: "USD", Balance: balance}
	mismatch, _ := NewMoney("EUR", decimal.NewFromInt(10))
	oversized, _ := NewMoney("USD", decimal.NewFromInt(60))

	err := wallet.Debit(mismatch)
	assert.ErrorIs(t, err, ErrCurrencyMismatchWallet)
	err = wallet.Debit(oversized)
	assert.ErrorIs(t, err, ErrInsufficientFunds)
}

func TestWalletCredit(t *testing.T) {
	balance, _ := NewMoney("USD", decimal.NewFromInt(20))
	wallet := Wallet{ID: "w1", Currency: "USD", Balance: balance}
	credit, _ := NewMoney("USD", decimal.NewFromInt(30))

	err := wallet.Credit(credit)
	assert.NoError(t, err)
	expected := decimal.NewFromInt(50)
	assert.True(t, wallet.Balance.Amount.Equal(expected))
}

func TestWalletCreditCurrencyMismatch(t *testing.T) {
	balance, _ := NewMoney("USD", decimal.NewFromInt(20))
	wallet := Wallet{ID: "w1", Currency: "USD", Balance: balance}
	credit, _ := NewMoney("EUR", decimal.NewFromInt(5))

	err := wallet.Credit(credit)
	assert.ErrorIs(t, err, ErrCurrencyMismatchWallet)
}
