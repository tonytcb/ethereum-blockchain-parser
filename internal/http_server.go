package main

import (
	"encoding/json"
	"fmt"
	"github.com/tonytcb/ethereum-blockchain-parser/pkg/domain"
	"io"
	"net/http"
)

type HTTPServer struct {
	port   string
	parser domain.Parser
}

func NewHTTPServer(port string, parser domain.Parser) *HTTPServer {
	return &HTTPServer{port: port, parser: parser}
}

func (s *HTTPServer) Start() error {
	http.HandleFunc("/subscribe", s.subscribeHandler)
	http.HandleFunc("/transactions", s.getTransactionsHandler)
	http.HandleFunc("/current-block", s.getCurrentBlockHandler)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", s.port), nil); err != nil {
			panic(err)
		}
	}()

	return nil
}

func (s *HTTPServer) Stop() error {
	return nil // @TODO implement me
}

func (s *HTTPServer) subscribeHandler(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = req.Body.Close()
	}()

	type payloadRequest struct {
		Address string `json:"address"`
	}

	data := &payloadRequest{}

	if err = json.Unmarshal(body, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if s.parser.Subscribe(data.Address) {
		w.WriteHeader(http.StatusCreated)
		return
	}

	http.Error(w, "error to subscribe", http.StatusServiceUnavailable)
}

func (s *HTTPServer) getTransactionsHandler(w http.ResponseWriter, req *http.Request) {
	var (
		address      = req.URL.Query().Get("address")
		transactions = s.parser.GetTransactions(address)
		output       = map[string]interface{}{
			"count":        len(transactions),
			"transactions": transactions,
		}
	)

	response, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(response)
}

func (s *HTTPServer) getCurrentBlockHandler(w http.ResponseWriter, req *http.Request) {
	// @TODO implement me
}
