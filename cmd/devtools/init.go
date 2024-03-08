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
	"strings"

	"github.com/spf13/cobra"
)

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
	// TODO: make the home dir configuratble
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
		return
	}
	defaultDir := filepath.Join(homeDir, ".astria")

	dataDir := "data"
	dataPath := filepath.Join(defaultDir, dataDir)
	createDir(dataPath)

	downloadDir := "local-dev-astria"
	fullPath := filepath.Join(defaultDir, downloadDir)

	fmt.Println("Local dev files placed in: ", fullPath)
	createDir(fullPath)
	recreateEnvFile(fullPath)
	recreateCometbftAndSequencerGenesisData(fullPath)

	for _, bin := range Binaries {

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

//go:embed config/local.env.example
var embeddedEnvironmentFile embed.FS

// TODO: add error handling
func recreateEnvFile(path string) {
	// Determine the user's home directory
	// TODO: replace homeDir with chose dir when custom home dir is implemented
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get user home directory: %v", err)
	}

	// Read the content from the embedded file
	data, err := fs.ReadFile(embeddedEnvironmentFile, "config/local.env.example")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
	}

	// Convert data to a string and replace "~" with the user's home directory
	content := strings.ReplaceAll(string(data), "~", homeDir)

	// Specify the path for the new file
	newPath := filepath.Join(path, ".env")

	// Create a new file
	newFile, err := os.Create(newPath)
	if err != nil {
		log.Fatalf("failed to create new file: %v", err)
	}
	defer newFile.Close()

	// Write the data to the new file
	_, err = newFile.WriteString(content)
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
