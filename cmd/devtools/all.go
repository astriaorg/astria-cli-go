/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// cleanallCmd represents the cleanall command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Clean all local data including binaries and config files.",
	Long:  "Clean all local data including binaries and config files. `dev init` will need to be run again to get the binaries and config files back.",
	Run: func(cmd *cobra.Command, args []string) {
		runCleanAll()
	},
}

func runCleanAll() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
		return
	}
	defaultDataDir := filepath.Join(homePath, ".astria")

	cleanCmd := exec.Command("rm", "-rf", defaultDataDir)
	if err := cleanCmd.Run(); err != nil {
		panic(err)
	}
}

func init() {
	cleanCmd.AddCommand(allCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
