package ethjsonrpc

const (
	defaultJSONRPCVersion = "2.0"

	ethNewFilterMethod        = "eth_newFilter"
	ethUninstallFilterMethod  = "eth_uninstallFilter"
	ethGetFilterChangesMethod = "eth_getFilterChanges"
)

type requestPayload struct {
	ID      int64       `json:"id"`
	JSONRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type newFilterParamsObjects struct {
	Address string `json:"address"`
}

type newFilterParams struct {
	Params []newFilterParamsObjects
}

type getFilterChangesParams struct {
	Addresses []string `json:"address"`
}

type uninstallFilterParams struct {
	Addresses []string `json:"address"`
}

func newRequestPayload(id int64, method string, params interface{}) requestPayload {
	return requestPayload{
		JSONRpc: defaultJSONRPCVersion,
		ID:      id,
		Method:  method,
		Params:  params,
	}
}
