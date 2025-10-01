package repositories

import (
	"context"
	"errors"
	"sync"

	"github.com/example/payment-gateway/internal/domain"
)

// InMemoryWalletRepository is a thread-safe wallet store for demos and tests.
type InMemoryWalletRepository struct {
	mu      sync.RWMutex
	wallets map[string]*domain.Wallet
}

// NewInMemoryWalletRepository creates the repository with preloaded wallets.
func NewInMemoryWalletRepository(initial map[string]*domain.Wallet) *InMemoryWalletRepository {
	store := make(map[string]*domain.Wallet, len(initial))
	for k, v := range initial {
		walletCopy := *v
		store[k] = &walletCopy
	}
	return &InMemoryWalletRepository{wallets: store}
}

// GetByID returns a wallet copy to avoid external mutation.
func (r *InMemoryWalletRepository) GetByID(_ context.Context, id string) (*domain.Wallet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	wallet, ok := r.wallets[id]
	if !ok {
		return nil, errors.New("wallet not found")
	}
	copy := *wallet
	return &copy, nil
}

// Save persists the wallet state.
func (r *InMemoryWalletRepository) Save(_ context.Context, wallet *domain.Wallet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := *wallet
	r.wallets[wallet.ID] = &copy
	return nil
}
