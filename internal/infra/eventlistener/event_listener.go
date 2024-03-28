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
	Add(context.Context, []domain.Transaction) error
}

type EthJSONAPI interface {
	NewFilter(ctx context.Context, address string) (string, error)
	FetchTransactions(ctx context.Context, address string) ([]domain.Transaction, error)
	RemoveFilter(ctx context.Context, address string) error
}

type EventListener struct {
	logger      *slog.Logger
	cfg         *config.Config
	api         EthJSONAPI
	storage     TransactionsStorage
	stopPooling chan struct{}
	mu          sync.Mutex
	filter      string
}

func (e *EventListener) Subscribe(ctx context.Context, address string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.filter != "" {
		return errors.New("Already subscribed")
	}

	filter, err := e.api.NewFilter(ctx, address)
	if err != nil {
		return errors.Wrap(err, "failed to create filter")
	}

	e.filter = filter

	go e.startPooling(filter)

	return nil
}

func (e *EventListener) startPooling(filter string) {
	ticker := time.NewTicker(e.cfg.PoolingTime)

	e.mu.Lock()
	e.stopPooling = make(chan struct{})
	e.mu.Unlock()

	defer e.logger.Info("Stopped pooling")

	for {
		select {
		case <-ticker.C:
			ctx := context.Background()

			transactions, err := e.api.FetchTransactions(ctx, filter)
			if err != nil {
				e.logger.Error("Failed to fetch transactions", "error", err)
				continue
			}

			if err = e.storage.Add(ctx, transactions); err != nil {
				e.logger.Error("Failed to store transactions", "error", err)
			}

		case <-e.stopPooling:
			ticker.Stop()
			return
		}
	}
}

func (e *EventListener) Unsubscribe() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.api.RemoveFilter(context.Background(), e.filter); err != nil {
		return errors.Wrap(err, "Failed to remove filter")
	}

	close(e.stopPooling)

	return nil
}
