package domain

import (
	"context"
	"errors"
	"log/slog"
	"math"
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []Transaction
}

type Repository interface {
	GetTransactions(ctx context.Context, address string) ([]Transaction, error)
	GetLatestBlock(ctx context.Context) (int64, error)
}

type EventListener interface {
	Listen(ctx context.Context, address string) error
}

type parser struct {
	logger        *slog.Logger
	repo          Repository
	eventListener EventListener
}

func NewParser(repo Repository, eventListener EventListener) Parser {
	return &parser{
		logger:        slog.Default(), // @TODO to make it optional via options
		repo:          repo,
		eventListener: eventListener,
	}
}

func (p *parser) GetCurrentBlock() int {
	blockNumber, err := p.repo.GetLatestBlock(context.Background())
	if err != nil {
		return 0
	}

	// Just to keep the proposed interface, we are returning zero in case of conversion overflow
	if blockNumber > math.MaxInt32 {
		return 0
	}

	return int(blockNumber)
}

func (p *parser) GetTransactions(address string) []Transaction {
	transactions, err := p.repo.GetTransactions(context.Background(), address)
	if err != nil && errors.Is(err, ErrAddressNotFound) {
		return []Transaction{}
	}

	return transactions
}

func (p *parser) Subscribe(address string) bool {
	// @TODO validate address

	err := p.eventListener.Listen(context.Background(), address)
	if err != nil {
		p.logger.Error("Failed to subscribe to address", "error", err)
	}

	return err == nil
}
