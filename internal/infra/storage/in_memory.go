package storage

import (
	"context"
	"sync"

	"github.com/tonytcb/ethereum-blockchain-parser/internal/domain"
)

type InMemory struct {
	mu           sync.RWMutex
	transactions map[string][]domain.Transaction
}

func (s *InMemory) Add(_ context.Context, address string, transactions []domain.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.transactions[address]; !ok {
		s.transactions[address] = make([]domain.Transaction, 0)
	}

	s.transactions[address] = append(s.transactions[address], transactions...)

	return nil
}

func (s *InMemory) GetTransactions(_ context.Context, address string) ([]domain.Transaction, error) {
	s.mu.RLock()
	defer s.mu.Unlock()

	if _, ok := s.transactions[address]; !ok {
		return []domain.Transaction{}, domain.ErrAddressNotFound
	}

	return s.transactions[address], nil
}
