package devtools

import (
	"archive/tar"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the local development environment.",
	Long:  `The init command will download the necessary binaries, create new directories for file organisation, and create an environment file for running a minimal Astria stack locally.`,
	Run:   runInitialization,
}

func init() {
	devCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("instance", "i", DefaultInstance, "Choose where the local-dev-astria directory will be created. Defaults to \"default\" if not provided.")
}

func runInitialization(c *cobra.Command, args []string) {
	// Get the instance name from the -i flag or use the default
	instance := c.Flag("instance").Value.String()
	err := IsInstanceNameValid(instance)
	if err != nil {
		log.WithError(err).Error("Error getting --instance flag")
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
		return
	}
	// TODO: make the default home dir configurable
	defaultDir := filepath.Join(homeDir, ".astria")
	instanceDir := filepath.Join(defaultDir, instance)

	fmt.Println("Creating new instance in:", instanceDir)

	// dataDir := "data"
	dataPath := filepath.Join(instanceDir, DataDir)
	createDir(dataPath)

	// downloadDir := "local-dev-astria"
	fullPath := filepath.Join(instanceDir, BinariesDir)
	fmt.Println("Local dev files placed in: ", fullPath)
	createDir(fullPath)

	recreateEnvFile(instanceDir, fullPath)
	recreateCometbftAndSequencerGenesisData(fullPath)

	for _, bin := range Binaries {
		downloadAndUnpack(bin.Url, fullPath, bin.Name)
	}

	initCometbft(instanceDir, DataDir, BinariesDir)

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
func recreateEnvFile(instancDir string, path string) {
	// Read the content from the embedded file
	data, err := fs.ReadFile(embeddedEnvironmentFile, "config/local.env.example")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
	}

	// Convert data to a string and replace "~" with the user's home directory
	content := strings.ReplaceAll(string(data), "~", instancDir)

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

// extractTarGz extracts a .tar.gz file to dest.
func extractTarGz(dest string, gzipStream io.Reader) error {
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
		target := filepath.Join(dest, header.Name)

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
	// TODO: should this be configuratble?
	err = os.Remove(dest)
	if err != nil {
		log.Fatalf("Failed to delete downloaded %s.tar.gz file: %v", packageName, err)
	}
	fmt.Printf("%s downloaded and extracted successfully.\n", packageName)
}

func initCometbft(defaultDir string, dataDirName string, binDirName string) {
	fmt.Println("Initializing CometBFT:")
	cometbftDataPath := filepath.Join(defaultDir, dataDirName, ".cometbft")

	// verify that cometbft was downloaded and extracted to the correct location
	cometbftBin := filepath.Join(defaultDir, binDirName, "cometbft")
	if !exists(cometbftBin) {
		fmt.Println("Error: cometbft binary not found here", cometbftBin)
		fmt.Println("\tCannot continue with initialization.")
		return
	}

	cometbftCmdPath := filepath.Join(defaultDir, binDirName, "cometbft")

	initCmdArgs := []string{"init", "--home", cometbftDataPath}
	initCmd := exec.Command(cometbftCmdPath, initCmdArgs...)

	fmt.Println("Running:", initCmd)

	_, err := initCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command", initCmd, ":", err)
		return
	} else {
		fmt.Println("\tSuccess")
	}

	// create the comand to replace the defualt genesis.json with the
	// configured one
	initGenesisJsonPath := filepath.Join(defaultDir, binDirName, "genesis.json")
	endGenesisJsonPath := filepath.Join(defaultDir, dataDirName, ".cometbft/config/genesis.json")
	copyArgs := []string{initGenesisJsonPath, endGenesisJsonPath}
	copyCmd := exec.Command("cp", copyArgs...)

	_, err = copyCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command", copyCmd, ":", err)
		return
	}
	fmt.Println("Copied genesis.json to", endGenesisJsonPath)

	// create the comand to replace the defualt priv_validator_key.json with the
	// configured one
	initPrivValidatorJsonPath := filepath.Join(defaultDir, binDirName, "priv_validator_key.json")
	endPrivValidatorJsonPath := filepath.Join(defaultDir, dataDirName, ".cometbft/config/priv_validator_key.json")
	copyArgs = []string{initPrivValidatorJsonPath, endPrivValidatorJsonPath}
	copyCmd = exec.Command("cp", copyArgs...)

	_, err = copyCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command", copyCmd, ":", err)
		return
	}
	fmt.Println("Copied priv_validator_key.json to", endPrivValidatorJsonPath)

	// update the cometbft config.toml file to have the proper block time
	cometbftConfigPath := filepath.Join(defaultDir, dataDirName, ".cometbft/config/config.toml")
	oldValue := `timeout_commit = "1s"`
	newValue := `timeout_commit = "2s"`

	if err := replaceInFile(cometbftConfigPath, oldValue, newValue); err != nil {
		fmt.Println("Error updating the file:", cometbftConfigPath, ":", err)
		return
	} else {
		fmt.Println("Successfully updated", cometbftConfigPath)
	}
}

// replaceInFile replaces oldValue with newValue in the file at filename.
// it is used here to update the block time in the cometbft config.toml file.
func replaceInFile(filename, oldValue, newValue string) error {
	// Read the original file.
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read the file: %w", err)
	}

	// Perform the replacement.
	modifiedContent := strings.ReplaceAll(string(content), oldValue, newValue)

	// Write the modified content to a new temporary file.
	tmpFilename := filename + ".tmp"
	if err := os.WriteFile(tmpFilename, []byte(modifiedContent), 0666); err != nil {
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}

	// Rename the original file to filename.bak.
	backupFilename := filename + ".bak"
	if err := os.Rename(filename, backupFilename); err != nil {
		return fmt.Errorf("failed to rename original file to backup: %w", err)
	}

	// Rename the temporary file to the original file name.
	if err := os.Rename(tmpFilename, filename); err != nil {
		// Attempt to restore the original file if renaming fails.
		os.Rename(backupFilename, filename)
		return fmt.Errorf("failed to rename temporary file to original: %w", err)
	}

	return nil
}
