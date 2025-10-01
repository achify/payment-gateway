package ports

import (
	"context"

	"github.com/example/payment-gateway/internal/domain"
)

// WalletRepository provides access to wallet persistence.
type WalletRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Wallet, error)
	Save(ctx context.Context, wallet *domain.Wallet) error
}
