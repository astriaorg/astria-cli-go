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
		runInitialization()
	},
}

func runInitialization() {
	// TODO: make the dir name configuratble
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting cwd:", err)
		return
	}

	dataDir := "data"
	dataPath := filepath.Join(cwd, dataDir)
	createDir(dataPath)

	downloadDir := "local-dev-astria"
	fullPath := filepath.Join(cwd, downloadDir)

	fmt.Println("Local dev files placed in: ", fullPath)
	createDir(fullPath)
	recreateEnvFile(fullPath)
	recreateCometbftAndSequencerGenesisData(fullPath)
	recreateMprocsFile(fullPath)
	recreateJustfile(fullPath)

	// TODO: make the binaries list configurable based on target os and arch
	binaries := []Binary{
		{"cometbft", "https://github.com/cometbft/cometbft/releases/download/v0.37.4/cometbft_0.37.4_darwin_arm64.tar.gz"},
		{"astria-sequencer", "https://github.com/astriaorg/astria/releases/download/sequencer-v0.9.0/astria-sequencer-aarch64-apple-darwin.tar.gz"},
		{"astria-composer", "https://github.com/astriaorg/astria/releases/download/composer-v0.4.0/astria-composer-aarch64-apple-darwin.tar.gz"},
		{"astria-conductor", "https://github.com/astriaorg/astria/releases/download/conductor-v0.12.0/astria-conductor-aarch64-apple-darwin.tar.gz"},
	}

	for _, bin := range binaries {

		downloadAndUnpack(bin.Url, fullPath, bin.Name)
	}

}

//go:embed config/genesis.json
var embeddedCometbftGenesisFile embed.FS

//go:embed config/priv_validator_key.json
var embeddedCometbftValidatorFile embed.FS

func recreateCometbftAndSequencerGenesisData(path string) {
	// Read the content from the embedded file
	genesisData, err := fs.ReadFile(embeddedCometbftGenesisFile, "config/genesis.json")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
	}
	// Read the content from the embedded file
	validatorData, err := fs.ReadFile(embeddedCometbftValidatorFile, "config/priv_validator_key.json")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
	}

	// Specify the path for the new file
	newGenesisPath := filepath.Join(path, "genesis.json")
	newValidatorPath := filepath.Join(path, "priv_validator_key.json")

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
	fmt.Println("Cometbft genesis data created successfully.")
	fmt.Println("Cometbft validator data created successfully.")
}

//go:embed config/justfile
var embeddedJustfile embed.FS

// TODO: add error handling
func recreateJustfile(path string) {
	// Read the content from the embedded file
	data, err := fs.ReadFile(embeddedJustfile, "config/justfile")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
	}

	// Specify the path for the new file
	newPath := filepath.Join(path, "justfile")

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
	fmt.Println("Justfile created successfully.")
}

//go:embed config/mprocs.yaml
var embeddedMprocsFile embed.FS

// TODO: add error handling
func recreateMprocsFile(path string) {
	// Read the content from the embedded file
	data, err := fs.ReadFile(embeddedMprocsFile, "config/mprocs.yaml")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
	}

	// Specify the path for the new file
	newPath := filepath.Join(path, "mprocs.yaml")

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
	fmt.Println("Mprocs file created successfully.")
}

//go:embed config/local.env.example
var embeddedEnvironmentFile embed.FS

// TODO: add error handling
func recreateEnvFile(path string) {
	// Read the content from the embedded file
	data, err := fs.ReadFile(embeddedEnvironmentFile, "config/local.env.example")
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
	fmt.Println("Local .env file created successfully.")
}

// TODO: add error handling
func createDir(dirName string) {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

}

// downloadFile downloads a file from the specified URL to the given local path.
func downloadFile(url, filepath string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// extractTarGz extracts a .tar.gz file to the current directory.
func extractTarGz(placePath string, gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(placePath, header.Name)

		// the following switch could also be done using if/else statements
		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			// handle normal file
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
}

// TODO: add error handling
func downloadAndUnpack(url string, placePath string, packageName string) {
	// Check if the file already exists
	if _, err := os.Stat(filepath.Join(placePath, packageName)); err == nil {
		fmt.Printf("%s already exists. Skipping download.\n", packageName)
		return
	}
	fmt.Printf("Downloading: (%s, %s)\n", packageName, url)

	// Download the file
	dest := filepath.Join(placePath, packageName+".tar.gz")
	if err := downloadFile(url, dest); err != nil {
		panic(err)
	}
	// Open the downloaded .tar.gz file
	file, err := os.Open(dest)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Extract the contents
	if err := extractTarGz(placePath, file); err != nil {
		panic(err)
	}

	// Delete the .tar.gz file
	err = os.Remove(dest)
	if err != nil {
		log.Fatalf("Failed to delete downloaded %s.tar.gz file: %v", packageName, err)
	}
	fmt.Printf("%s downloaded and extracted successfully.\n", packageName)
}

func init() {
	devCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// TODO: add a "path" flag to the init command
	// initCmd.Flags().StringP("path", "p", "", "Choose where the local-dev-astria directory will be created. Defaults to the current working directory.")
}
