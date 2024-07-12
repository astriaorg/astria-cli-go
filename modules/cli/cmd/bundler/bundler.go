package bundler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	primproto "buf.build/gen/go/astria/primitives/protocolbuffers/go/astria/primitive/v1"
	txproto "buf.build/gen/go/astria/protocol-apis/protocolbuffers/go/astria/protocol/transactions/v1alpha1"
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// generateCmd represents the balances command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates and prints a sample JSON from a protobuf.",
	Run:   GenerateJSONFromProtos,
}

func GenerateJSONFromProtos(c *cobra.Command, args []string) {
	toAddr := &primproto.Address{
		Bech32M: "sampleToAddress",
	}
	amount1, _ := convertToUint128("69")
	amount2, _ := convertToUint128("99999999999")
	amount3, _ := convertToUint128("1234567890")

	tx := &txproto.UnsignedTransaction{
		Params: &txproto.TransactionParams{
			ChainId: "sampleChainId",
			// Nonce:   1,
		},
		Actions: []*txproto.Action{
			{
				Value: &txproto.Action_TransferAction{
					TransferAction: &txproto.TransferAction{
						To:       toAddr,
						Amount:   amount1,
						Asset:    "sampleAsset",
						FeeAsset: "sampleFeeAsset",
					},
				},
			},
			{
				Value: &txproto.Action_TransferAction{
					TransferAction: &txproto.TransferAction{
						To:       toAddr,
						Amount:   amount2,
						Asset:    "sampleAsset",
						FeeAsset: "sampleFeeAsset",
					},
				},
			},
			{
				Value: &txproto.Action_TransferAction{
					TransferAction: &txproto.TransferAction{
						To:       toAddr,
						Amount:   amount3,
						Asset:    "sampleAsset",
						FeeAsset: "sampleFeeAsset",
					},
				},
			},
		},
	}

	// Encode to JSON
	marshaller := jsonpb.Marshaler{}
	jsonStr, err := marshaller.MarshalToString(tx)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		os.Exit(1)
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(jsonStr), "", "  "); err != nil {
		fmt.Println("Error pretty converting JSON:", err)
		os.Exit(1)
	}

	fmt.Println(prettyJSON.String())

}

// readCmd represents the balances command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Reads a sample JSON from a protobuf.",
	Run:   readCmdHandler,
}

func readCmdHandler(c *cobra.Command, args []string) {

	str := `{
		"actions": [
		  {
			"transferAction": {
			  "to": {
				"bech32m": "sampleToAddress"
			  },
			  "amount": {
				"lo": "69"
			  },
			  "asset": "sampleAsset",
			  "feeAsset": "sampleFeeAsset"
			}
		  },
		  {
			"transferAction": {
			  "to": {
				"bech32m": "sampleToAddress"
			  },
			  "amount": {
				"lo": "99999999999"
			  },
			  "asset": "sampleAsset",
			  "feeAsset": "sampleFeeAsset"
			}
		  },
		  {
			"transferAction": {
			  "to": {
				"bech32m": "sampleToAddress"
			  },
			  "amount": {
				"lo": "1234567890"
			  },
			  "asset": "sampleAsset",
			  "feeAsset": "sampleFeeAsset"
			}
		  }
		],
		"params": {
		  "chainId": "sampleChainId"
		}
	  }`

	person := &txproto.UnsignedTransaction{}

	// Unmarshal the JSON string into the protobuf message
	if err := protojson.Unmarshal([]byte(str), person); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// can then set the nonce after reading
	person.Params.Nonce = 1
	// Print the protobuf message
	fmt.Println(person)
}

func init() {
	cmd.RootCmd.AddCommand(generateCmd)
	cmd.RootCmd.AddCommand(readCmd)
}
