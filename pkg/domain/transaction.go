package domain

import "github.com/pkg/errors"

type Transaction struct {
	Hash               string `json:"transactionHash"`
	Address            string `json:"address"`
	BlockNumber        string `json:"blockNumber"`
	BlockHash          string `json:"blockHash"`
	DecimalBlockNumber int64  `json:"decimalBlockNumber"`
}

func NewTransaction(hash, address, blockNumber, blockHash string) (Transaction, error) {
	decimalBlockNumber, err := hexToDecimalString(blockNumber)
	if err != nil {
		return Transaction{}, errors.Wrap(err, "invalid block number")
	}

	return Transaction{
		Hash:               hash,
		Address:            address,
		BlockNumber:        blockNumber,
		BlockHash:          blockHash,
		DecimalBlockNumber: decimalBlockNumber,
	}, nil
}
