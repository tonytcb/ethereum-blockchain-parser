package ethjsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/tonytcb/ethereum-blockchain-parser/pkg/domain"
)

type Options func(*EthJSONRpc)

func WithHTTPClient(c *http.Client) Options {
	return func(e *EthJSONRpc) {
		e.httpClient = c
	}
}

type Config struct {
	APIURL         string
	RequestTimeout time.Duration
}

type EthJSONRpc struct {
	cfg        *Config
	httpClient *http.Client
	currentID  int64
}

func NewEthJSONRpc(cfg *Config, opts ...Options) *EthJSONRpc {
	e := &EthJSONRpc{
		currentID: 0,
		cfg:       cfg,
		httpClient: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func (e *EthJSONRpc) NewFilter(ctx context.Context, address string) (string, error) {
	e.currentID++

	payload := newRequestPayload(e.currentID, ethNewFilterMethod, []struct {
		Address string `json:"address"`
	}{
		{Address: address},
	})

	resPayload, err := e.doPost(ctx, payload)
	if err != nil {
		return "", errors.Wrap(err, "error reading response body")
	}

	var response newFilterResponse

	if err = json.Unmarshal(resPayload, &response); err != nil {
		return "", errors.Wrap(err, "error unmarshalling response")
	}

	if response.Error != nil {
		return "", errors.Errorf("error response: %s", response.Error.Message)
	}

	return response.Result, nil
}

func (e *EthJSONRpc) FetchTransactions(ctx context.Context, filter string) ([]domain.Transaction, error) {
	e.currentID++

	payload := newRequestPayload(e.currentID, ethGetFilterChangesMethod, []string{filter})

	resPayload, err := e.doPost(ctx, payload)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	var response getFilterChangesResponse

	if err = json.Unmarshal(resPayload, &response); err != nil {
		return nil, errors.Wrap(err, "error unmarshalling response")
	}

	if response.Error != nil {
		return nil, errors.Errorf("error response: %s", response.Error.Message)
	}

	var transactions = make([]domain.Transaction, len(response.Result))

	for i, v := range response.Result {
		t, err := domain.NewTransaction(
			v.TransactionHash,
			v.Address,
			v.BlockNumber,
			v.BlockHash,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error creating new transaction")
		}

		transactions[i] = t
	}

	return transactions, nil
}

func (e *EthJSONRpc) RemoveFilter(ctx context.Context, address string) error {
	e.currentID++

	payload := newRequestPayload(e.currentID, ethUninstallFilterMethod, []string{address})

	resPayload, err := e.doPost(ctx, payload)
	if err != nil {
		return errors.Wrap(err, "error reading response body")
	}

	var response uninstallFilterResponse

	if err = json.Unmarshal(resPayload, &response); err != nil {
		return errors.Wrap(err, "error unmarshalling response")
	}

	if response.Error != nil {
		return errors.Errorf("error response: %s", response.Error.Message)
	}

	return nil
}

func (e *EthJSONRpc) doPost(ctx context.Context, payload requestPayload) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling payload")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.cfg.APIURL, bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "error creating new request")
	}

	req.Header.Set("Content-Type", "application/json")

	const maxRetries = 10

	var res *http.Response

	// TODO Add maxRetries and sleep time on configuration

	for i := 1; i <= maxRetries; i++ {
		res, err = e.httpClient.Do(req)
		if err != nil {
			continue
		}
		if res.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(time.Millisecond * 100 * time.Duration(i))
	}

	if err != nil {
		return nil, errors.Wrap(err, "error executing request")
	}
	defer func() {
		_ = res.Body.Close()
	}()

	resPayload, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	return resPayload, nil
}
