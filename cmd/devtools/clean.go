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

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans the local development environment data.",
	Long:  `Cleans the local development environment data. Does not remove the binaries or the config files.`,
	Run: func(cmd *cobra.Command, args []string) {
		runClean()
	},
}

func runClean() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
		return
	}
	defaultDataDir := filepath.Join(homePath, ".astria/data")

	cleanCmd := exec.Command("rm", "-rf", defaultDataDir)
	if err := cleanCmd.Run(); err != nil {
		panic(err)
	}

	err = os.MkdirAll(defaultDataDir, 0755) // Read & execute by everyone, write by owner.
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
}

func init() {
	devCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
