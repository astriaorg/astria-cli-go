# go-sequencer-client

## Usage

```go
package main

import (
  "context"
  "crypto/sha256"
  "fmt"

  sqproto "buf.build/gen/go/astria/astria/protocolbuffers/go/astria/sequencer/v1alpha1"
  client "github.com/astriaorg/astria-cli-go/modules/go-sequencer-client/client"
)

func main() {
  signer, err := client.GenerateSigner()
  if err != nil {
    panic(err)
  }

  // default tendermint RPC endpoint
  c, err := client.NewClient("http://localhost:26657")
  if err != nil {
    panic(err)
  }

  rollupId := sha256.Sum256([]byte("test-chain"))
  tx := &sqproto.UnsignedTransaction{
    Nonce: 1,
    Actions: []*sqproto.Action{
      {
        Value: &sqproto.Action_SequenceAction{
          SequenceAction: &sqproto.SequenceAction{
            RollupId: rollupId[:],
            Data:     []byte("test-data"),
          },
        },
      },
    },
  }

  signed, err := signer.SignTransaction(tx)
  if err != nil {
    panic(err)
  }

  resp, err := c.BroadcastTxSync(context.Background(), signed)
  if err != nil {
    panic(err)
  }

  fmt.Println(resp)
}
```
