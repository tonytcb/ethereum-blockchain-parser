package domain

import (
	"reflect"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	type args struct {
		hash        string
		address     string
		blockNumber string
		blockHash   string
	}
	tests := []struct {
		name    string
		args    args
		want    Transaction
		wantErr bool
	}{
		{
			name: "should return expected transaction",
			args: args{
				hash:        "0x3585a53d05d27d622e3c5be8c34b5c4a6ad101929f3efbbac84dfd0129c67cc7",
				address:     "0x35fa164735182de50811e8e2e824cfb9b6118ac2",
				blockNumber: "0x12a30a0",
				blockHash:   "0xacfdbd8d63fbb7cdea95b2d85be04ee5ecd8eef222a2861e246b8454bcf3952c",
			},
			want: Transaction{
				Hash:               "0x3585a53d05d27d622e3c5be8c34b5c4a6ad101929f3efbbac84dfd0129c67cc7",
				Address:            "0x35fa164735182de50811e8e2e824cfb9b6118ac2",
				BlockNumber:        "0x12a30a0",
				BlockHash:          "0xacfdbd8d63fbb7cdea95b2d85be04ee5ecd8eef222a2861e246b8454bcf3952c",
				DecimalBlockNumber: 19542176,
			},
			wantErr: false,
		},
		{
			name: "should error on invalid block number",
			args: args{
				blockNumber: "0x",
			},
			want:    Transaction{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTransaction(tt.args.hash, tt.args.address, tt.args.blockNumber, tt.args.blockHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}
