package domain

import "strings"

// CurrencyPair represents a tradable pair of currencies.
type CurrencyPair struct {
	Base  string
	Quote string
}

// NewCurrencyPair builds a normalised currency pair.
func NewCurrencyPair(base, quote string) CurrencyPair {
	return CurrencyPair{Base: strings.ToUpper(base), Quote: strings.ToUpper(quote)}
}

// IsValid returns whether the currency pair contains two different currencies.
func (p CurrencyPair) IsValid() bool {
	return p.Base != "" && p.Quote != "" && p.Base != p.Quote
}
