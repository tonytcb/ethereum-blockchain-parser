# ethereum-blockchain-parser

Implements a `parser` to allow to query transactions for subscribed addresses.

The `parser` makes use of the [Ethereum JSON RPC Interface](https://ethereum.org/en/developers/docs/apis/json-rpc) to create an event listener, processing new transactions created for the subscribed address on the blockchain.

## Design

This package was designed to be easily extensible, by simply implementing its Storage or Event Listener interface, or pointing the Ethereum JSON RPC API to a different provider, as well.

## Parser Interface

```go
type Parser interface {
    GetCurrentBlock() int
    Subscribe(address string) bool
    GetTransactions(address string) []Transaction
}
```

## How to use

```go

import (
    "github.com/tonytcb/ethereum-blockchain-parser/pkg/domain"
    "github.com/tonytcb/ethereum-blockchain-parser/pkg/infra/ethjsonrpc"
    "github.com/tonytcb/ethereum-blockchain-parser/pkg/infra/eventlistener"
    "github.com/tonytcb/ethereum-blockchain-parser/pkg/infra/storage"
)

func main() {
    repository = storage.NewInMemory()
    api        = ethjsonrpc.NewEthJSONRpc(&ethjsonrpc.Config{
        APIURL:         cfg.EthereumRPCAPIURL,
        RequestTimeout: cfg.RequestTimeout,
    })
    eventListener = eventlistener.NewPoolingEventListener(
        ctx,
        api,
        repository,
        eventlistener.WithLogger(logger),
        eventlistener.WithConfig(&eventlistener.Config{PoolingTime: cfg.PoolingTime}),
    )
    parser = domain.NewParser(repository, eventListener)
	    
    err := parser.Subscribe("0x1234567890")
    transactions, _ := parser.GetTransactions("0x1234567890")
    blockNumber := parser.GetCurrentBlock()
}

```

## Example

As a simple example of usage, was implemented an application consuming this package and exposing an HTTP API at [internal package](./internal).
