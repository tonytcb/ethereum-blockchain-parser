package main

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"

	"github.com/tonytcb/ethereum-blockchain-parser/pkg/domain"
	"github.com/tonytcb/ethereum-blockchain-parser/pkg/infra/ethjsonrpc"
	"github.com/tonytcb/ethereum-blockchain-parser/pkg/infra/eventlistener"
	"github.com/tonytcb/ethereum-blockchain-parser/pkg/infra/storage"
)

type Application struct {
	cfg           *Config
	logger        *slog.Logger
	eventListener domain.EventListener
	httpServer    *HTTPServer
}

func NewApplication(_ context.Context, cfg *Config, logger *slog.Logger) *Application {
	var (
		repository = storage.NewInMemory()
		api        = ethjsonrpc.NewEthJSONRpc(&ethjsonrpc.Config{
			APIURL:         cfg.EthereumRPCAPIURL,
			RequestTimeout: cfg.RequestTimeout,
		})
		eventListener = eventlistener.NewPoolingEventListener(
			api,
			repository,
			eventlistener.WithLogger(logger),
			eventlistener.WithConfig(&eventlistener.Config{PoolingTime: cfg.PoolingTime}),
		)

		parser = domain.NewParser(repository, eventListener)

		httpServer = NewHTTPServer(cfg.HTTPPort, parser)
	)

	return &Application{
		logger:        logger,
		cfg:           cfg,
		eventListener: eventListener,
		httpServer:    httpServer,
	}
}

//nolint:unparam
func (a *Application) Run(ctx context.Context) error {
	errGroup, _ := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		return a.httpServer.Start()
	})

	return nil
}

func (a *Application) Stop() error {
	return a.httpServer.Stop()
}
