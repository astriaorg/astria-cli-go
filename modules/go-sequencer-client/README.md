# go-sequencer-client

The `go-sequencer-client` is a Go package that enables interacting with an
Astria sequencer.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Example](#example)

## Installation

To use the `go-sequencer-client`, install the following packages:

1. Sequencer Client:

   ```bash
   go get github.com/astriaorg/astria-cli-go/modules/go-sequencer-client
   ```

2. Protobuf Types:
   The Astria sequencer uses ProtoBuf for message passing. A full list of the
   APIs and Primitives can be found [here](https://buf.build/astria).

   ```bash
   go get buf.build/gen/go/astria/primitives/protocolbuffers/go
   go get buf.build/gen/go/astria/protocol-apis/protocolbuffers/go
   ```

3. Bech32m Package:
   The Astria sequencer uses "astria" prefixed `bech32m` addresses.

   ```bash
   go get github.com/astriaorg/astria-cli-go/modules/bech32m
   ```

## Usage

All sequencer client methods can be found in the [client.go](./client/client.go)
file.

## Example

The following example demonstrates how to create a sequencer client, then build
and send a transaction to a sequencer:

```go
package main

import (
  "context"
  "crypto/ed25519"
  "encoding/hex"
  "fmt"

  "github.com/astriaorg/astria-cli-go/modules/bech32m"
  "github.com/astriaorg/astria-cli-go/modules/go-sequencer-client/client"
  txproto "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"
  primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
)

// create a sequencer client and send a transfer
func main() {
  // create a new sequencer client pointing to the default CometBFT RPC endpoint
  c, err := client.NewClient("http://localhost:26657")
 if err != nil {
  log.WithError(err).Error("Error creating sequencer client")
  panic(err)
 }

  // parse a private key into an ed25519.PrivateKey
 privKeyBytes, err := hex.DecodeString("hex private key string")
 if err != nil {
  panic(err)
 }
 from := ed25519.NewKeyFromSeed(privKeyBytes)

  // create a FROM address from the private key
 signer := client.NewSigner(from)
 fromAddr := signer.Address()
 addr, err := bech32m.EncodeFromBytes("astria", fromAddr)
 if err != nil {
  log.WithError(err).Error("Failed to encode address")
  panic(err)
 }

  // automatically get the nonce for the FROM account
 nonce, err := c.GetNonce(ctx, addr.String())
 if err != nil {
  log.WithError(err).Error("Error getting nonce")
  panic(err)
 }

  // convert the transfer amount to a uint128
  bigInt := new(big.Int)
 _, ok := bigInt.SetString("1000", 10) // the transfer amount is 1000
 if !ok {
  return nil, fmt.Errorf("failed to convert string to big.Int")
 }
 if bigInt.Sign() < 0 {
  panic(fmt.Errorf("negative number not allowed"))
 } else if bigInt.BitLen() > 128 {
  panic(fmt.Errorf("value overflows Uint128"))
 }
 lo := bigInt.Uint64()
 hi := bigInt.Rsh(bigInt, 64).Uint64()
 transferAmount := &primproto.Uint128{
  Lo: lo,
  Hi: hi,
 }

  // build the transaction
 unsignedTx := &txproto.UnsignedTransaction{
  Params: &txproto.TransactionParams{
   ChainId: "test-sequencer",
   Nonce:   nonce,
  },
  Actions: []*txproto.Action{
   {
    Value: &txproto.Action_TransferAction{
     TransferAction: &txproto.TransferAction{
      To:       "astria prefixed bech32m TO address",
      Amount:   transferAmount,
      Asset:    "asset",
      FeeAsset: "feeAsset",
     },
    },
   },
  },
 }

  signedTx, err := signer.SignTransaction(unsignedTx)
  if err != nil {
    panic(err)
  }

  sendAsync := false
  resp, err := c.BroadcastTx(context.Background(), signedTx, sendAsync)
  if err != nil {
    panic(err)
  }

  fmt.Println(resp)
}
```
