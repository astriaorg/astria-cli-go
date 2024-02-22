/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"archive/tar"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Binary struct {
	Name string
	Url  string
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the local development environment.",
	Long:  `The init command will download the nessesary binaries, create new directories for file organisation, and create an environment file for running a minimal Astria stack locally.`,
	Run: func(cmd *cobra.Command, args []string) {
		initDev()
	},
}

func initDev() {
	// TODO: make the dir name configuratble
	downloadDir := "local-dev-astria"
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fullPath := filepath.Join(cwd, downloadDir)

	createDevDir(fullPath)
	recreateEnvFile(fullPath)
	recreateCometbftAndSequencerGenesisData(fullPath)

	binaries := []Binary{
		{"cometbft", "https://github.com/cometbft/cometbft/releases/download/v0.37.4/cometbft_0.37.4_darwin_arm64.tar.gz"},
		{"astria-sequencer", "https://github.com/astriaorg/astria/releases/download/sequencer-v0.9.0/astria-sequencer-aarch64-apple-darwin.tar.gz"},
		{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v0.4.0/astria-composer-aarch64-apple-darwin.tar.gz"},
		{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v0.12.0/astria-conductor-aarch64-apple-darwin.tar.gz"},
	}

	// fileURL := "https://example.com/file.tar.gz"
	// filePath := "downloaded_file.tar.gz"

	for _, bin := range binaries {
		fmt.Printf("Downloading: (%s, %s)\n", bin.Name, bin.Url)

		downloadPath := filepath.Join(fullPath, bin.Name)
		if err := DownloadAndExtractFile(downloadPath, bin.Url); err != nil {
			log.Fatal("Download and extraction failed: ", err)
		}

		log.Println("File downloaded and extracted successfully.")
	}

}

//go:embed genesis.json
var embeddedCometbftGenesisFile embed.FS

//go:embed priv_validator_key.json
var embeddedCometbftValidatorFile embed.FS

func recreateCometbftAndSequencerGenesisData(path string) {
	genesisPath := "genesis.json"
	validatorPath := "priv_validator_key.json"
	// Read the content from the embedded file
	genesisData, err := fs.ReadFile(embeddedCometbftGenesisFile, genesisPath)
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
	}
	// Read the content from the embedded file
	validatorData, err := fs.ReadFile(embeddedCometbftValidatorFile, validatorPath)
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
	}

	// Specify the path for the new file
	newGenesisPath := filepath.Join(path, genesisPath)
	newValidatorPath := filepath.Join(path, validatorPath)

	// Create a new file
	newGenesisFile, err := os.Create(newGenesisPath)
	if err != nil {
		log.Fatalf("failed to create new file: %v", err)
	}
	defer newGenesisFile.Close()
	newValidatorFile, err := os.Create(newValidatorPath)
	if err != nil {
		log.Fatalf("failed to create new file: %v", err)
	}
	defer newValidatorFile.Close()

	// Write the data to the new file
	_, err = newGenesisFile.Write(genesisData)
	if err != nil {
		log.Fatalf("failed to write data to new file: %v", err)
	}
	_, err = newValidatorFile.Write(validatorData)
	if err != nil {
		log.Fatalf("failed to write data to new file: %v", err)
	}
}

//go:embed local.env.example
var embeddedEnvironmentFile embed.FS

// TODO: add error handling
func recreateEnvFile(path string) {
	// Read the content from the embedded file
	data, err := fs.ReadFile(embeddedEnvironmentFile, "local.env.example")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
	}

	// Specify the path for the new file
	newPath := filepath.Join(path, ".env")

	// Create a new file
	newFile, err := os.Create(newPath)
	if err != nil {
		log.Fatalf("failed to create new file: %v", err)
	}
	defer newFile.Close()

	// Write the data to the new file
	_, err = newFile.Write(data)
	if err != nil {
		log.Fatalf("failed to write data to new file: %v", err)
	}
}

// TODO: add error handling
func createDevDir(dirName string) {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

}

func DownloadAndExtractFile(path string, url string) error {
	// Download section (same as before)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Extraction section
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	// Extract tar contents
	for {
		header, err := tr.Next()

		// If no more files are found return
		if err == io.EOF {
			break
		}

		// Return any other error
		if err != nil {
			return err
		}

		// Check the type of the file
		target := filepath.Join(path, "extracted", header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			println("dir")
			// Create directory
			if err := os.MkdirAll(target, 0777); err != nil {
				return err
			}
		case tar.TypeReg:
			println("file")

			// Create file
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			println("file")

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			println("file")

			outFile.Close()
			println("file")

		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
