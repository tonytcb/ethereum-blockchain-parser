package storage

import (
	"context"
	"testing"

	"github.com/tonytcb/ethereum-blockchain-parser/pkg/domain"
)

func TestNewInMemory(t *testing.T) {
	var (
		t1 = domain.Transaction{BlockNumber: "1", DecimalBlockNumber: 1}
		t2 = domain.Transaction{BlockNumber: "2", DecimalBlockNumber: 2}
		t3 = domain.Transaction{BlockNumber: "3", DecimalBlockNumber: 3}

		ctx = context.Background()
	)
	storage := NewInMemory()

	_, err := storage.GetLatestBlock(ctx)
	if err != domain.ErrBlockNotFound {
		t.Fatalf("expected error due to address not found")
	}

	err = storage.Add(ctx, "1", []domain.Transaction{t1, t2, t3})
	if err != nil {
		t.Fatalf("expected error to be nil")
	}

	all, err := storage.GetTransactions(ctx, "1")
	if err != nil {
		t.Fatalf("expected error to be nil")
	}

	if len(all) != 3 {
		t.Fatalf("expected 3 transactions")
	}

	err = storage.UpdateLastBlock(ctx, "1", 1)
	if err != nil {
		t.Fatalf("expected error to be nil")
	}

	err = storage.UpdateLastBlock(ctx, "1", 3)
	if err != nil {
		t.Fatalf("expected error to be nil")
	}

	latestBlock, err := storage.GetLatestBlock(ctx)
	if err != nil {
		t.Fatalf("expected error to be nil")
	}
	if latestBlock != 3 {
		t.Fatalf("expected latest block to be 3, got: %v", latestBlock)
	}
}
