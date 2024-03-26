package app

import (
	"context"

	"github.com/tonytcb/ethereum-blockchain-parser/internal/app/config"
)

type Application struct {
	cfg *config.Config
}

func NewApplication(ctx context.Context, cfg *config.Config) (*Application, error) {
	return &Application{
		cfg: cfg,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	return nil
}

func (a *Application) Stop() error {
	return nil
}
