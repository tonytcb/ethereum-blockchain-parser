package domain

type Transaction struct {
	Hash        string `json:"transactionHash"`
	Address     string `json:"address"`
	BlockNumber string `json:"blockNumber"`
	BlockHash   string `json:"blockHash"`
}
