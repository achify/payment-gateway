package domain

import "github.com/shopspring/decimal"

// TransferQuote encapsulates the outcome of calculating a transfer between two currencies.
type TransferQuote struct {
	Pair          CurrencyPair
	ExchangeRate  decimal.Decimal
	SendAmount    Money
	ReceiveAmount Money
	Fee           Money
	TotalDebit    Money
	Allowed       bool
}
