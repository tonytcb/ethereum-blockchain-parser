//tags: integration

package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/tonytcb/ethereum-blockchain-parser/pkg/domain"
	"github.com/tonytcb/ethereum-blockchain-parser/pkg/infra/eventlistener"
	"github.com/tonytcb/ethereum-blockchain-parser/pkg/infra/storage"
	"github.com/tonytcb/ethereum-blockchain-parser/pkg/mocks"
)

func TestParserIntegrationTest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		api         = mocks.NewEthJSONAPI(t)
		poolingTime = time.Millisecond * 10
		address     = "0x35fA164735182de50811E8e2E824cFb9B6118ac2"
		filterID    = "0x1"
	)

	api.EXPECT().NewFilter(mock.Anything, address).Return(filterID, nil).Once()
	api.EXPECT().RemoveFilter(mock.Anything, filterID).Return(nil).Maybe()

	var count = 0
	api.EXPECT().FetchTransactions(mock.Anything, filterID).RunAndReturn(func(ctx context.Context, s string) ([]domain.Transaction, error) {
		count++

		switch count {
		case 1:
			return []domain.Transaction{}, nil

		case 2:
			return []domain.Transaction{
				{
					Hash:               "0x1",
					Address:            address,
					BlockNumber:        "0x1",
					BlockHash:          "0x1",
					DecimalBlockNumber: 1,
				},
			}, nil

		case 3:
			return []domain.Transaction{
				{
					Hash:               "0x2",
					Address:            address,
					BlockNumber:        "0x1",
					BlockHash:          "0x1",
					DecimalBlockNumber: 1,
				},
				{
					Hash:               "0x3",
					Address:            address,
					BlockNumber:        "0x2",
					BlockHash:          "0x2",
					DecimalBlockNumber: 2,
				},
			}, nil
		default:
			return []domain.Transaction{}, nil
		}
	}).Maybe()

	var (
		repository    = storage.NewInMemory()
		eventListener = eventlistener.NewPoolingEventListener(
			ctx,
			api,
			repository,
			eventlistener.WithConfig(&eventlistener.Config{PoolingTime: poolingTime}),
		)
		parser = domain.NewParser(repository, eventListener)
	)

	t.Run("GetTransactions should return zero when not subscribed", func(t *testing.T) {
		transactions := parser.GetTransactions(address)
		if len(transactions) != 0 {
			t.Errorf("expected zero transactions, got %d", len(transactions))
		}

		lastBlock := parser.GetCurrentBlock()
		if lastBlock != 0 {
			t.Errorf("expected last block to be 0, got %d", lastBlock)
		}
	})

	t.Run("Should subscribe and return 3 transactions on GetTransactions successfully", func(t *testing.T) {
		if success := parser.Subscribe(address); !success {
			t.Fatalf("expected true, got false")
		}

		// waiting for the event listener to fetch the transactions
		time.Sleep(poolingTime * 5)

		transactions := parser.GetTransactions(address)
		if len(transactions) != 3 {
			t.Errorf("expected zero transactions, got %d", len(transactions))
		}

		lastBlock := parser.GetCurrentBlock()
		if lastBlock != 2 {
			t.Errorf("expected last block to be 2, got %d", lastBlock)
		}

		if err := eventListener.Unsubscribe(ctx, address); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}
