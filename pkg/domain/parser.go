package domain

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"regexp"
)

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []Transaction
}

type RepositoryReader interface {
	GetTransactions(ctx context.Context, address string) ([]Transaction, error)
	GetLatestBlock(ctx context.Context) (int64, error)
}

type EventListener interface {
	Listen(ctx context.Context, address string) error
}

type Options func(*parser)

func WithLogger(l *slog.Logger) Options {
	return func(e *parser) {
		e.logger = l
	}
}

type parser struct {
	logger        *slog.Logger
	repo          RepositoryReader
	eventListener EventListener
}

func NewParser(repo RepositoryReader, eventListener EventListener, opts ...Options) Parser {
	p := &parser{
		logger:        slog.Default(),
		repo:          repo,
		eventListener: eventListener,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *parser) GetCurrentBlock() int {
	blockNumber, err := p.repo.GetLatestBlock(context.Background())
	if err != nil {
		return 0
	}

	// Just to keep the proposed interface, we are returning zero in case of conversion overflow
	if blockNumber > math.MaxInt32 {
		p.logger.Error("Block number int32 overflow", "blockNumber", blockNumber)
		return 0
	}

	return int(blockNumber)
}

func (p *parser) GetTransactions(address string) []Transaction {
	transactions, err := p.repo.GetTransactions(context.Background(), address)
	if (err != nil && errors.Is(err, ErrAddressNotFound)) || (transactions == nil) {
		return []Transaction{}
	}

	return transactions
}

func (p *parser) Subscribe(address string) bool {
	if !isValidAddress(address) {
		return false
	}

	err := p.eventListener.Listen(context.Background(), address)
	if err != nil {
		p.logger.Error("Failed to subscribe to address", "error", err)
	}

	return err == nil
}

// isValidAddress checks if the given string is a valid Ethereum address
// @see https://goethereumbook.org/en/address-check/
func isValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}
