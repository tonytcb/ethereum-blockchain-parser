package main

import (
	"context"
)

type Application struct {
	cfg *Config
}

func NewApplication(_ context.Context, cfg *Config) (*Application, error) {
	return &Application{
		cfg: cfg,
	}, nil
}

func (a *Application) Run(_ context.Context) error {
	return nil
}

func (a *Application) Stop() error {
	return nil
}
