package domain

import (
	"context"
	"log/slog"
)

//type Storage interface {
//	GetLastParsedBlock(blockID int) (int, error)
//}

// ============================ interfaces

type Parser interface {
	// GetCurrentBlock return the last parsed block
	GetCurrentBlock() int

	// Subscribe subscribes to an address
	Subscribe(address string) bool

	// GetTransactions returns a list of transactions based on address
	GetTransactions(address string) []Transaction
}

type EventListener interface {
	Subscribe(ctx context.Context, address string) error
}

// ============================ real implementations

type parser struct {
	logger        *slog.Logger
	eventListener EventListener
}

func (p *parser) GetCurrentBlock() int {
	//TODO implement me
	panic("implement me")
}

func (p *parser) Subscribe(address string) bool {
	// @TODO validate address

	err := p.eventListener.Subscribe(context.Background(), address)
	if err != nil {
		p.logger.Error("Failed to subscribe to address", "error", err)
	}

	return err == nil
}

func (p *parser) GetTransactions(address string) []Transaction {
	//TODO implement me
	panic("implement me")
}

// 1. creates an event listener to pool changes on new blocks
// 1.1 new filter on api
// 1.2 http pool based on filter id
// 2. query transactions based on new event
// 3. parse transactions