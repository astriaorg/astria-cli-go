//go:build integration_tests
// +build integration_tests

package integration_tests

import (
	"context"
	"testing"

	"github.com/astriaorg/astria-cli-go/modules/go-sequencer-client/client"
	"github.com/stretchr/testify/require"
)

func TestGetBalance(t *testing.T) {
	c, err := client.NewClient("http://localhost:26657")
	require.NoError(t, err)

	balance, err := c.GetBalances(context.Background(), "astria1hj8pc8vwcvrr7wswjulemzls4cm9mj5w5858df")
	require.NoError(t, err)
	require.Empty(t, balance)
}

func TestGetNonce(t *testing.T) {
	c, err := client.NewClient("http://localhost:26657")
	require.NoError(t, err)

	nonce, err := c.GetNonce(context.Background(), "astria1hj8pc8vwcvrr7wswjulemzls4cm9mj5w5858df")
	require.NoError(t, err)
	require.Equal(t, nonce, uint32(0))
}
