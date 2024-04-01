package domain_test

import (
	"bytes"
	"log/slog"
	"math"
	"reflect"
	"testing"

	"github.com/pkg/errors"

	"github.com/tonytcb/ethereum-blockchain-parser/pkg/domain"

	"github.com/stretchr/testify/mock"
	"github.com/tonytcb/ethereum-blockchain-parser/pkg/mocks"
)

func Test_parser_GetCurrentBlock(t *testing.T) {
	var logger = slog.New(slog.NewJSONHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))

	type fields struct {
		logger        *slog.Logger
		repo          func(*testing.T) domain.RepositoryReader
		eventListener domain.EventListener
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "should return 0 when there's no blocks",
			fields: fields{
				logger: logger,
				repo: func(t *testing.T) domain.RepositoryReader {
					repo := mocks.NewRepositoryReader(t)
					repo.EXPECT().GetLatestBlock(mock.Anything).Return(int64(0), domain.ErrBlockNotFound).Once()
					return repo
				},
				eventListener: nil,
			},
			want: 0,
		},
		{
			name: "should return 0 when repository's response overflow int32 max value",
			fields: fields{
				logger: logger,
				repo: func(t *testing.T) domain.RepositoryReader {
					repo := mocks.NewRepositoryReader(t)
					repo.EXPECT().GetLatestBlock(mock.Anything).Return(math.MaxInt32+1, nil).Once()
					return repo
				},
				eventListener: nil,
			},
			want: 0,
		},
		{
			name: "should return 65 on success",
			fields: fields{
				logger: logger,
				repo: func(t *testing.T) domain.RepositoryReader {
					repo := mocks.NewRepositoryReader(t)
					repo.EXPECT().GetLatestBlock(mock.Anything).Return(int64(65), nil).Once()
					return repo
				},
				eventListener: nil,
			},
			want: 65,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := domain.NewParser(
				tt.fields.repo(t),
				tt.fields.eventListener,
			)

			if got := p.GetCurrentBlock(); got != tt.want {
				t.Errorf("GetCurrentBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_GetTransactions(t *testing.T) {
	var (
		logger       = slog.New(slog.NewJSONHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
		transactions = []domain.Transaction{
			{
				Hash:    "0x1",
				Address: "0x123",
			},
		}
	)

	type fields struct {
		logger *slog.Logger
		repo   func(*testing.T) domain.RepositoryReader
	}
	type args struct {
		address string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []domain.Transaction
	}{
		{
			name: "should return empty transactions list on address not error",
			fields: fields{
				logger: logger,
				repo: func(t *testing.T) domain.RepositoryReader {
					repo := mocks.NewRepositoryReader(t)
					repo.EXPECT().GetTransactions(mock.Anything, "0x123").Return(nil, domain.ErrAddressNotFound).Once()
					return repo
				},
			},
			args: args{
				address: "0x123",
			},
			want: []domain.Transaction{},
		},
		{
			name: "should return empty transactions list on address any other error",
			fields: fields{
				logger: logger,
				repo: func(t *testing.T) domain.RepositoryReader {
					repo := mocks.NewRepositoryReader(t)
					repo.EXPECT().GetTransactions(mock.Anything, "0x123").Return(nil, errors.New("network error")).Once()
					return repo
				},
			},
			args: args{
				address: "0x123",
			},
			want: []domain.Transaction{},
		},
		{
			name: "should return 1 transaction on success",
			fields: fields{
				logger: logger,
				repo: func(t *testing.T) domain.RepositoryReader {
					repo := mocks.NewRepositoryReader(t)
					repo.EXPECT().GetTransactions(mock.Anything, "0x123").Return(transactions, nil).Once()
					return repo
				},
			},
			args: args{
				address: "0x123",
			},
			want: transactions,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := domain.NewParser(tt.fields.repo(t), nil)

			got := p.GetTransactions(tt.args.address)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTransactions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_Subscribe(t *testing.T) {
	var logger = slog.New(slog.NewJSONHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))

	type fields struct {
		logger        *slog.Logger
		eventListener func(*testing.T) domain.EventListener
	}
	type args struct {
		address string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "should return true on success",
			fields: fields{
				logger: logger,
				eventListener: func(t *testing.T) domain.EventListener {
					el := mocks.NewEventListener(t)
					el.EXPECT().Listen(mock.Anything, "0x35fA164735182de50811E8e2E824cFb9B6118ac2").Return(nil).Once()
					return el
				},
			},
			args: args{
				address: "0x35fA164735182de50811E8e2E824cFb9B6118ac2",
			},
			want: true,
		},
		{
			name: "should return false on invalid address",
			fields: fields{
				logger: logger,
				eventListener: func(*testing.T) domain.EventListener {
					return nil
				},
			},
			args: args{
				address: "0x35fA164735182de50811E8e2E824cFb9B6118ac2_",
			},
			want: false,
		},
		{
			name: "should return false when already subscribed",
			fields: fields{
				logger: logger,
				eventListener: func(t *testing.T) domain.EventListener {
					el := mocks.NewEventListener(t)
					el.EXPECT().Listen(mock.Anything, "0x35fA164735182de50811E8e2E824cFb9B6118ac2").
						Return(domain.ErrAlreadySubscribed).
						Once()
					return el
				},
			},
			args: args{
				address: "0x35fA164735182de50811E8e2E824cFb9B6118ac2",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := domain.NewParser(nil, tt.fields.eventListener(t), domain.WithLogger(tt.fields.logger))

			if got := p.Subscribe(tt.args.address); got != tt.want {
				t.Errorf("Subscribe() = %v, want %v", got, tt.want)
			}
		})
	}
}
