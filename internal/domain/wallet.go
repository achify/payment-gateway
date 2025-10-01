package domain

import (
	"errors"
)

var (
	ErrCurrencyMismatchWallet = errors.New("wallet currency mismatch")
	ErrInsufficientFunds      = errors.New("insufficient funds")
)

// Wallet represents a customer wallet with a single currency balance.
type Wallet struct {
	ID       string
	Currency string
	Balance  Money
}

// Debit deducts the provided amount from the wallet balance.
func (w *Wallet) Debit(amount Money) error {
	if w.Currency != amount.Currency {
		return ErrCurrencyMismatchWallet
	}
	if w.Balance.Amount.LessThan(amount.Amount) {
		return ErrInsufficientFunds
	}
	newBalance := w.Balance.Amount.Sub(amount.Amount)
	w.Balance.Amount = newBalance
	return nil
}

// Credit adds the provided amount to the wallet balance.
func (w *Wallet) Credit(amount Money) error {
	if w.Currency != amount.Currency {
		return ErrCurrencyMismatchWallet
	}
	w.Balance.Amount = w.Balance.Amount.Add(amount.Amount)
	return nil
}
