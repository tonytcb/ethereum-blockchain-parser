package eventlistener

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/tonytcb/ethereum-blockchain-parser/internal/app/config"
	"github.com/tonytcb/ethereum-blockchain-parser/internal/domain"
)

type TransactionsStorage interface {
	Add(ctx context.Context, address string, transactions []domain.Transaction) error
}

type EthJSONAPI interface {
	NewFilter(ctx context.Context, address string) (string, error)
	FetchTransactions(ctx context.Context, address string) ([]domain.Transaction, error)
	RemoveFilter(ctx context.Context, address string) error
}

type PoolingEventListener struct {
	mu          sync.Mutex
	logger      *slog.Logger
	cfg         *config.Config // @TODO replace by a local config struct
	api         EthJSONAPI
	storage     TransactionsStorage
	stopPooling map[string]chan struct{}
	filters     map[string]string
}

func NewPoolingEventListener(
	logger *slog.Logger,
	cfg *config.Config,
	api EthJSONAPI,
	storage TransactionsStorage,
) *PoolingEventListener {
	return &PoolingEventListener{
		logger:      logger,
		cfg:         cfg,
		api:         api,
		storage:     storage,
		stopPooling: make(map[string]chan struct{}),
		filters:     make(map[string]string),
	}
}

func (e *PoolingEventListener) Subscribe(ctx context.Context, address string) error {
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

	defer e.logger.Info("Stopped pooling address %s", address)

	for {
		select {
		case <-ticker.C:
			e.logger.Debug("Pooling transactions against address %s, filter %s", address, filter)

			ctx := context.Background()

			transactions, err := e.api.FetchTransactions(ctx, filter)
			if err != nil {
				e.logger.Error("Failed to fetch transactions", "error", err)
				continue
			}

			if err = e.storage.Add(ctx, address, transactions); err != nil {
				e.logger.Error("Failed to store transactions", "error", err)
			}

		case <-stopPoolingCh:
			ticker.Stop()
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
