package eventlistener

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/tonytcb/ethereum-blockchain-parser/pkg/domain"
)

const (
	defaultPoolingTime = 1 * time.Second
)

type Repository interface {
	Add(ctx context.Context, address string, transactions []domain.Transaction) error
	UpdateLastBlock(ctx context.Context, address string, blockNumber int64) error
}

type EthJSONAPI interface {
	NewFilter(ctx context.Context, address string) (string, error)
	FetchTransactions(ctx context.Context, address string) ([]domain.Transaction, error)
	RemoveFilter(ctx context.Context, address string) error
}

type Options func(*PoolingEventListener)

func WithConfig(cfg *Config) Options {
	return func(e *PoolingEventListener) {
		e.cfg = cfg
	}
}

func WithLogger(l *slog.Logger) Options {
	return func(e *PoolingEventListener) {
		e.logger = l
	}
}

type Config struct {
	PoolingTime time.Duration
}

type PoolingEventListener struct {
	mu          sync.Mutex
	logger      *slog.Logger
	cfg         *Config
	api         EthJSONAPI
	repository  Repository
	stopPooling map[string]chan struct{}
	filters     map[string]string
}

func NewPoolingEventListener(
	api EthJSONAPI,
	storage Repository,
	opts ...Options,
) *PoolingEventListener {
	e := &PoolingEventListener{
		logger:      slog.Default(),
		cfg:         &Config{PoolingTime: defaultPoolingTime},
		api:         api,
		repository:  storage,
		stopPooling: make(map[string]chan struct{}),
		filters:     make(map[string]string),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func (e *PoolingEventListener) Listen(ctx context.Context, address string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.filters[address]; ok {
		return domain.ErrAlreadySubscribed
	}

	filter, err := e.api.NewFilter(ctx, address)
	if err != nil {
		return errors.Wrap(err, "failed to create filter")
	}

	e.filters[address] = filter

	go e.startPooling(address, filter)

	return nil
}

func (e *PoolingEventListener) startPooling(address string, filter string) {
	ticker := time.NewTicker(e.cfg.PoolingTime)
	stopPoolingCh := make(chan struct{})

	e.mu.Lock()
	e.stopPooling[address] = stopPoolingCh
	e.mu.Unlock()

	for {
		select {
		case <-ticker.C:
			e.logger.Debug("Pooling transactions", "address", address, "filter", filter)

			ctx := context.Background()

			transactions, err := e.api.FetchTransactions(ctx, filter)
			if err != nil {
				e.logger.Error("Failed to fetch transactions", "error", err)
				continue
			}

			if len(transactions) == 0 {
				continue
			}

			if err = e.repository.Add(ctx, address, transactions); err != nil {
				e.logger.Error("Failed to store transactions", "error", err)
				continue
			}

			lastBlock := e.highestBlockNumber(transactions)
			if err = e.repository.UpdateLastBlock(ctx, address, lastBlock); err != nil {
				e.logger.Error("Failed to update last block", "error", err)
			}

		case <-stopPoolingCh:
			if err := e.api.RemoveFilter(context.Background(), filter); err != nil {
				e.logger.Error("Failed to remove filter", "error", err)
			}

			ticker.Stop()

			e.logger.Info("Stopped pooling", "address", address)

			return
		}
	}
}

func (e *PoolingEventListener) Unsubscribe(ctx context.Context, address string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	filter, ok := e.filters[address]
	if !ok {
		return domain.ErrNotSubscribed
	}

	if err := e.api.RemoveFilter(ctx, filter); err != nil {
		return errors.Wrap(err, "Failed to remove filter")
	}

	close(e.stopPooling[address])

	return nil
}

func (e *PoolingEventListener) highestBlockNumber(transactions []domain.Transaction) int64 {
	length := len(transactions)

	if length == 0 {
		return 0
	}

	highest := transactions[0].DecimalBlockNumber

	for i := 0; i < length; i++ {
		blockNumber := transactions[i].DecimalBlockNumber
		if transactions[i].DecimalBlockNumber > highest {
			highest = blockNumber
		}
	}

	return highest
}
