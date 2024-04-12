package sequencer

import (
	"encoding/json"
	"strconv"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:    "transfer [amount] [to] --privkey=[privkey]",
	Short:  "Transfer tokens from one account to another.",
	Args:   cobra.ExactArgs(2),
	PreRun: cmd.SetLogLevel,
	Run:    runTransfer,
}

func init() {
	sequencerCmd.AddCommand(transferCmd)

	transferCmd.Flags().String("privkey", "", "The private key of the account from which to transfer tokens.")
	transferCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer.")
	transferCmd.Flags().Bool("json", false, "Output in JSON format.")

	err := transferCmd.MarkFlagRequired("privkey")
	if err != nil {
		log.WithError(err).Error("Error marking flag as required")
		panic(err)
	}
}

func runTransfer(cmd *cobra.Command, args []string) {
	printJSON := cmd.Flag("json").Value.String() == "true"

	amount := args[0]
	to := args[1]

	url := cmd.Flag("url").Value.String()
	from := cmd.Flag("privkey").Value.String()

	opts := sequencer.TransferOpts{
		SequencerURL: url,
		FromKey:      from,
		ToAddress:    to,
		Amount:       amount,
	}
	res, err := sequencer.Transfer(opts)
	if err != nil {
		log.WithError(err).Error("Error transferring tokens")
		panic(err)
	}

	if printJSON {
		j, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.WithError(err).Error("Error marshalling to JSON")
			panic(err)
		}
		pterm.Println(string(j))
		return
	} else {
		header := []string{"From", "To", "Nonce", "Amount", "TxHash"}
		var rows [][]string
		rows = append(rows, []string{res.From, res.To, strconv.Itoa(int(res.Nonce)), res.Amount, res.TxHash})
		data := append([][]string{header}, rows...)
		output, err := pterm.DefaultTable.WithHasHeader().WithSeparator(" ").WithData(data).Srender()
		if err != nil {
			log.WithError(err).Error("Error rendering table")
			panic(err)
		}
		pterm.Println(output)
	}
}
