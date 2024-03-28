package sequencer

import (
	"fmt"

	"github.com/astria/astria-cli-go/internal/sequencer"
	"github.com/spf13/cobra"
)

// getBalanceCmd represents the getBalance command
var getBalanceCmd = &cobra.Command{
	Use:   "get-balance",
	Short: "Retrieves and prints the balance of an account.",
	Long: `Use this command to retrieve and print the balance of an account.

Usage: astria-cli-go sequencer get-balance <address> --url <url>`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:  runGetBalance,
}

const DefaultSequencerURL = "http://127.0.0.1:26657"

func init() {
	sequencerCmd.AddCommand(getBalanceCmd)
	getBalanceCmd.Flags().String("url", DefaultSequencerURL, "The URL of the sequencer to retrieve the balance from.")
}

func runGetBalance(cmd *cobra.Command, args []string) {
	// TODO - arg validation?
	address := args[0]
	url := cmd.Flag("url").Value.String()

	balance, err := sequencer.GetBalance(address, url)
	if err != nil {
		fmt.Println("Error getting balance:", err)
		return
	}

	fmt.Println(balance)
}
