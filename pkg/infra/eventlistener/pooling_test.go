package eventlistener

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/tonytcb/ethereum-blockchain-parser/pkg/domain"

	"github.com/stretchr/testify/mock"
	"github.com/tonytcb/ethereum-blockchain-parser/pkg/mocks"
)

func TestPoolingEventListener_Listen(t *testing.T) {
	var logger = slog.New(slog.NewJSONHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))

	var transactions = []domain.Transaction{
		{
			Hash:               "0x4a1353904f4a2ef953a2c9db7f770a2c4ad50f7273ef1bfd0cd6b62a924859ca",
			Address:            "0x35fa164735182de50811e8e2e824cfb9b6118ac2",
			BlockNumber:        "0x12a3193",
			BlockHash:          "0xf04376490c0ffc1cf3cefecc5a36b4eb24e76909b28453f3e8592f22fadd972b",
			DecimalBlockNumber: 19542419,
		},
		{
			Hash:               "0x3585a53d05d27d622e3c5be8c34b5c4a6ad101929f3efbbac84dfd0129c67cc7",
			Address:            "0x35fa164735182de50811e8e2e824cfb9b6118ac2",
			BlockNumber:        "0x12a30a0",
			BlockHash:          "0xacfdbd8d63fbb7cdea95b2d85be04ee5ecd8eef222a2861e246b8454bcf3952c",
			DecimalBlockNumber: 19542176,
		},
	}

	type fields struct {
		ctx    context.Context
		logger *slog.Logger
		cfg    *Config
		api    func(*testing.T) EthJSONAPI
		repo   func(*testing.T) RepositoryWriter
	}
	type args struct {
		ctx     context.Context
		address string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		waitTime time.Duration
		wantErr  bool
	}{
		{
			name: "should start listening for new transactions but process none",
			fields: fields{
				ctx:    context.Background(),
				logger: logger,
				cfg:    &Config{PoolingTime: time.Millisecond * 10},
				api: func(t *testing.T) EthJSONAPI {
					api := mocks.NewEthJSONAPI(t)
					api.EXPECT().NewFilter(mock.Anything, "0x123").Return("0x1", nil).Once()
					api.EXPECT().FetchTransactions(mock.Anything, "0x1").Return([]domain.Transaction{}, nil).Once()
					api.EXPECT().RemoveFilter(mock.Anything, "0x1").Return(nil).Maybe()
					return api
				},
				repo: func(t *testing.T) RepositoryWriter {
					return mocks.NewRepositoryWriter(t)
				},
			},
			args: args{
				ctx:     context.Background(),
				address: "0x123",
			},
			waitTime: time.Millisecond * 15,
			wantErr:  false,
		},
		{
			name: "should start listening for new transactions and process two transactions",
			fields: fields{
				ctx:    context.Background(),
				logger: logger,
				cfg:    &Config{PoolingTime: time.Millisecond * 50},
				api: func(t *testing.T) EthJSONAPI {
					api := mocks.NewEthJSONAPI(t)
					api.EXPECT().NewFilter(mock.Anything, "0x123").Return("0x2", nil).Once()
					api.EXPECT().FetchTransactions(mock.Anything, "0x2").Return(transactions, nil).Once()
					api.EXPECT().RemoveFilter(mock.Anything, "0x2").Return(nil).Once()
					return api
				},
				repo: func(t *testing.T) RepositoryWriter {
					repo := mocks.NewRepositoryWriter(t)
					repo.EXPECT().Add(mock.Anything, "0x123", transactions).Return(nil).Once()
					repo.EXPECT().UpdateLastBlock(mock.Anything, "0x123", int64(19542419)).Return(nil).Once()
					return repo
				},
			},
			args: args{
				ctx:     context.Background(),
				address: "0x123",
			},
			waitTime: time.Millisecond * 60,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			e := NewPoolingEventListener(
				tt.fields.ctx,
				tt.fields.api(t),
				tt.fields.repo(t),
				WithLogger(tt.fields.logger),
				WithConfig(tt.fields.cfg),
			)

			if err := e.Listen(tt.args.ctx, tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("Listen() error = %v, wantErr %v", err, tt.wantErr)
			}

			time.Sleep(tt.waitTime) // sleep time to process async events

			err := e.Unsubscribe(tt.args.ctx, tt.args.address)
			if err != nil {
				t.Errorf("Unsubscribe() error = %v", err)
			}
		})
	}
}
