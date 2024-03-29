package ethjsonrpc

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type newFilterResponse struct {
	ID     int            `json:"id"`
	Result string         `json:"result"`
	Error  *errorResponse `json:"error"`
}

type getFilterChangesResponse struct {
	ID     int `json:"id"`
	Result []struct {
		Address          string   `json:"address"`
		Topics           []string `json:"topics"`
		Data             string   `json:"data"`
		BlockNumber      string   `json:"blockNumber"`
		TransactionHash  string   `json:"transactionHash"`
		TransactionIndex string   `json:"transactionIndex"`
		BlockHash        string   `json:"blockHash"`
		LogIndex         string   `json:"logIndex"`
		Removed          bool     `json:"removed"`
	} `json:"result"`
	Error *errorResponse `json:"error"`
}

type uninstallFilterResponse struct {
	ID    int            `json:"id"`
	Error *errorResponse `json:"error"`
}
