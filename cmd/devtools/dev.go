package devtools

import (
	"fmt"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "TODO: short description",
	Long:  `TODO: longer description`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: print out the help for all dev commands
		fmt.Println("dev called")
	},
}

// TODO: add a func here to print out an explanation for how to use the dev command

func init() {
	cmd.RootCmd.AddCommand(devCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
