package repositories

import (
	"fmt"
	"sync"

	"github.com/itsXanadium/singa-riset-backend/domain"
)

type MemoryRepository struct {
	verifications map[string]domain.Verification
	mu            sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{verifications: make(map[string]domain.Verification)}
}
func (r *MemoryRepository) Store(email string, v domain.Verification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.verifications[email] = v
	return nil
}
func (r *MemoryRepository) Fetch(email string) (domain.Verification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	v, exists := r.verifications[email]
	if !exists {
		return domain.Verification{}, fmt.Errorf("Verification Not Found")
	}
	return v, nil
}

func (r *MemoryRepository) Delete(email string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.verifications, email)
	return nil
}
