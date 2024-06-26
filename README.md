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
    var (
        logger     = slog.Default()
        repository = storage.NewInMemory()
        api        = ethjsonrpc.NewEthJSONRpc(&ethjsonrpc.Config{
            APIURL:         "https://cloudflare-eth.com",
            RequestTimeout: 3 * time.Millisecond,
        })
        eventListener = eventlistener.NewPoolingEventListener(
            context.Background(),
            api,
            repository,
            eventlistener.WithLogger(logger),
            eventlistener.WithConfig(&eventlistener.Config{PoolingTime: 1 * time.Second}),
        )

        parser = domain.NewParser(repository, eventListener)
    )
	    
    err := parser.Subscribe("0x1234567890")
    transactions, err := parser.GetTransactions("0x1234567890")
    blockNumber := parser.GetCurrentBlock()
}

```

## Example

As a simple example of usage, was implemented an application consuming this package and exposing an HTTP API at [internal package](./internal).

## Links and references

- [Etherscan](https://etherscan.io)
- [Ethereum JSON RPC Interface](https://ethereum.org/en/developers/docs/apis/json-rpc)
- [All That Node](https://www.allthatnode.com/ethereum.dsrv)
- [Cloudflare ETH API](https://developers.cloudflare.com/web3/ethereum-gateway/reference/supported-api-methods/)
