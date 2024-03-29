package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/tonytcb/ethereum-blockchain-parser/pkg/domain"
)

type InMemory struct {
	mu           sync.RWMutex
	lastBlock    map[string]int64
	transactions map[string][]domain.Transaction
}

func NewInMemory() *InMemory {
	return &InMemory{
		mu:           sync.RWMutex{},
		lastBlock:    make(map[string]int64),
		transactions: make(map[string][]domain.Transaction),
	}
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

func (s *InMemory) UpdateLastBlock(_ context.Context, address string, blockNumber int64) error {
	if blockNumber <= 0 {
		return errors.New("block number must be greater than zero")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastBlock[address] = blockNumber

	return nil
}

func (s *InMemory) GetTransactions(_ context.Context, address string) ([]domain.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.transactions[address]; !ok {
		return []domain.Transaction{}, domain.ErrAddressNotFound
	}

	return s.transactions[address], nil
}

func (s *InMemory) GetLatestBlock(_ context.Context) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// returns first element, if present
	for _, block := range s.lastBlock {
		return block, nil
	}

	return 0, domain.ErrBlockNotFound
}
