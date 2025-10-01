package domain

// PayoutDestination contains information required to deliver funds to a recipient.
type PayoutDestination struct {
	Type          string
	AccountNumber string
	BankCode      string
	BankName      string
	Country       string
	Currency      string
	RecipientName string
}
