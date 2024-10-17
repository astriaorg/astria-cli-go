package client

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	accountsproto "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/accounts/v1"
	txproto "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transaction/v1"
	"github.com/cometbft/cometbft/rpc/client"
	"github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// BalanceResponse describes the response from a balance query.
// Should this live here?
type BalanceResponse struct {
	Denom   string   `json:"denom,omitempty"`
	Balance *big.Int `json:"balance,omitempty"`
}

// Client is an HTTP tendermint client.
type Client struct {
	websocket string
	client    *http.HTTP
}

func NewClient(url string) (*Client, error) {
	// Compile the regular expression
	re := regexp.MustCompile(`^[^:]+://`)

	c, err := http.New(url, "")
	if err != nil {
		return nil, err
	}

	// Replace and print results
	websocket := re.ReplaceAllString(url, "")
	websocket = "tcp://" + websocket
	return &Client{
		websocket: websocket,
		client:    c,
	}, nil
}

// BroadcastTx broadcasts a transaction. If async is true, the function will
// return immediately. The response seen is the generated data used for
// submitting the transaction. It does not confirm that the data has been
// included on chain. If async is false, the function will wait for the
// transaction to be seen on the network.
func (c *Client) BroadcastTx(ctx context.Context, tx *txproto.Transaction, async bool) (*coretypes.ResultBroadcastTx, error) {
	if async {
		return c.BroadcastTxAsync(ctx, tx)
	}
	return c.BroadcastTxSync(ctx, tx)
}

// BroadcastTxAsync broadcasts a transaction and returns immediately.
func (c *Client) BroadcastTxAsync(ctx context.Context, tx *txproto.Transaction) (*coretypes.ResultBroadcastTx, error) {
	bytes, err := proto.Marshal(tx)
	if err != nil {
		return nil, err
	}

	return c.client.BroadcastTxAsync(ctx, bytes)
}

// BroadcastTxSync broadcasts a transaction and waits for the response that
// confirms the transaction was included.
func (c *Client) BroadcastTxSync(ctx context.Context, tx *txproto.Transaction) (*coretypes.ResultBroadcastTx, error) {
	bytes, err := proto.Marshal(tx)
	if err != nil {
		return nil, err
	}
	result, resultErr := c.client.BroadcastTxSync(ctx, bytes)
	if resultErr != nil {
		return result, resultErr
	}
	// must check result.Code because cometbft doesn't return an error on tx failure
	if result.Code != 0 {
		return result, errors.New(result.Log)
	}

	// wait for the tx to be included in a block
	retryCount := 50
	retryInterval := 250 * time.Millisecond
	for i := 0; i < retryCount; i++ {
		t, err := c.client.Tx(ctx, result.Hash, true)
		if err != nil {
			// manually ignore error that says tx not found
			if strings.HasPrefix(err.Error(), "error in json rpc client, with http response metadata: (Status: 200 OK, Protocol HTTP/1.1). RPC error -32603") {
				log.Debug("tx not found, retrying...")
				// wait a short amount of time between retries
				time.Sleep(retryInterval)
				continue
			} else {
				return result, err
			}

		}
		if t != nil {
			return result, nil
		}
	}

	return result, fmt.Errorf("tx %s not found after %d retries", result.Hash, retryCount)
}

func (c *Client) GetBalances(ctx context.Context, addr string) ([]*BalanceResponse, error) {
	query := "accounts/balance/" + addr
	resp, err := c.client.ABCIQueryWithOptions(ctx, query, []byte{}, client.ABCIQueryOptions{
		Height: 0,
		Prove:  false,
	})
	if err != nil {
		return nil, err
	}

	if resp.Response.Code != 0 {
		return nil, errors.New(resp.Response.Log)
	}

	protoBalanceResp := &accountsproto.BalanceResponse{}
	err = proto.Unmarshal(resp.Response.Value, protoBalanceResp)
	if err != nil {
		return nil, err
	}

	return balanceResponseFromProto(protoBalanceResp), nil
}

func (c *Client) GetNonce(ctx context.Context, addr string) (uint32, error) {
	query := "accounts/nonce/" + addr
	resp, err := c.client.ABCIQueryWithOptions(ctx, query, []byte{}, client.ABCIQueryOptions{
		Height: 0,
		Prove:  false,
	})
	if err != nil {
		return 0, err
	}

	if resp.Response.Code != 0 {
		return 0, errors.New(resp.Response.Log)
	}

	nonceResp := &accountsproto.NonceResponse{}
	err = proto.Unmarshal(resp.Response.Value, nonceResp)
	if err != nil {
		return 0, err
	}

	return nonceResp.Nonce, nil
}

func (c *Client) GetBlock(ctx context.Context, height *int64) (*coretypes.ResultBlock, error) {
	block, err := c.client.Block(ctx, height)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (c *Client) GetBlockHeight(ctx context.Context) (int64, error) {
	block, err := c.client.Block(ctx, nil)
	if err != nil {
		return 0, err
	}

	return block.Block.Height, nil
}

func balanceResponseFromProto(resp *accountsproto.BalanceResponse) []*BalanceResponse {
	var balanceResponses []*BalanceResponse
	for _, balance := range resp.Balances {
		balanceResponses = append(balanceResponses, &BalanceResponse{
			Balance: ProtoU128ToBigInt(balance.Balance),
			Denom:   balance.Denom,
		})
	}
	return balanceResponses
}
