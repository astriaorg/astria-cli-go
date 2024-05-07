package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	util "github.com/astria/astria-cli-go/cmd/devtools/utilities"

	log "github.com/sirupsen/logrus"
)

// IsInstanceNameValidOrPanic checks if the instance name is valid and panics if it's not.
func IsInstanceNameValidOrPanic(instance string) {
	re, err := regexp.Compile(`^[a-z]+[a-z0-9]*(-[a-z0-9]+)*$`)
	if err != nil {
		log.WithError(err).Error("Error compiling regex")
		panic(err)
	}
	if !re.MatchString(instance) {
		log.Errorf("Invalid instance name: %s", instance)
		err := fmt.Errorf("invalid instance name: '%s'. Instance names must be lowercase, alphanumeric, and may contain dashes. It can't begin or end with a dash. No repeating dashes", instance)
		panic(err)
	}
}

// IsSequencerChainIdValidOrPanic checks if the instance name is valid and panics if it's not.
func IsSequencerChainIdValidOrPanic(id string) {
	if len(id) < 1 || len(id) > 50 {
		log.Errorf("Invalid sequencer chain id: %s", id)
		err := fmt.Errorf("invalid sequencer chain id: '%s'. The ChainId length must be within the range [1,50], can contain lowercase and uppercase letter, numerical digits, and the characters '-', '_', and '.'", id)
		panic(err)
	}

	re, err := regexp.Compile(`^[a-zA-Z0-9\-_\.]+$`)
	if err != nil {
		log.WithError(err).Error("Error compiling regex")
		panic(err)
	}
	if !re.MatchString(id) {
		log.Errorf("Invalid sequencer chain id: %s", id)
		err := fmt.Errorf("invalid sequencer chain id: '%s'. The ChainId length must be within the range [1,50], can contain lowercase and uppercase letter, numerical digits, and the characters '-', '_', and '.'", id)
		panic(err)
	}
}

//go:embed local.env.example
var embeddedLocalEnvironmentFile embed.FS

// RecreateLocalEnvFile creates a new local .env file at the specified path.
func RecreateLocalEnvFile(instanceDir string, path string) {
	// Read the content from the embedded file
	data, err := fs.ReadFile(embeddedLocalEnvironmentFile, "local.env.example")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
		panic(err)
	}

	// Convert data to a string and replace "~" with the user's home directory
	content := strings.ReplaceAll(string(data), "~", instanceDir)

	// Specify the path for the new file
	newPath := filepath.Join(path, ".env")

	// check if the local .env file already exists
	_, err = os.Stat(newPath)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", newPath)
		return
	}

	// Create a new file
	newFile, err := os.Create(newPath)
	if err != nil {
		log.Fatalf("failed to create new file: %v", err)
		panic(err)
	}
	defer newFile.Close()

	// Write the data to the new file
	_, err = newFile.WriteString(content)
	if err != nil {
		log.Fatalf("failed to write data to new file: %v", err)
		panic(err)
	}
	log.Infof("Local .env file created successfully: %s\n", newPath)
}

//go:embed remote.env.example
var embeddedRemoteEnvironmentFile embed.FS

// RecreateRemoteEnvFile creates a new remote .env file at the specified path.
func RecreateRemoteEnvFile(instanceDir string, path string) {
	// Read the content from the embedded file
	data, err := fs.ReadFile(embeddedRemoteEnvironmentFile, "remote.env.example")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
		panic(err)
	}

	// Specify the path for the new file
	newPath := filepath.Join(path, ".env")

	_, err = os.Stat(newPath)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", newPath)
		return
	}

	// Create a new file
	newFile, err := os.Create(newPath)
	if err != nil {
		log.Fatalf("failed to create new file: %v", err)
		panic(err)
	}
	defer newFile.Close()

	// Write the data to the new file
	_, err = newFile.WriteString(string(data))
	if err != nil {
		log.Fatalf("failed to write data to new file: %v", err)
		panic(err)
	}
	log.Infof("Remote .env file created successfully: %s\n", newPath)

}

//go:embed genesis.json
var embeddedCometbftGenesisFile embed.FS

//go:embed priv_validator_key.json
var embeddedCometbftValidatorFile embed.FS

// RecreateCometbftAndSequencerGenesisData creates a new CometBFT genesis.json and priv_validator_key.json file at the specified path.
func RecreateCometbftAndSequencerGenesisData(path, localNetworkName string) {
	// Read the content from the embedded file
	genesisData, err := fs.ReadFile(embeddedCometbftGenesisFile, "genesis.json")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
		panic(err)
	}
	// Unmarshal JSON into a map to update sequencer chain id
	var data map[string]interface{}
	if err := json.Unmarshal(genesisData, &data); err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}
	// update chain id and convert back to bytes
	data["chain_id"] = localNetworkName
	genesisData, err = json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling updated data to JSON: %s", err)
	}

	// Read the content from the embedded file
	validatorData, err := fs.ReadFile(embeddedCometbftValidatorFile, "priv_validator_key.json")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
		panic(err)
	}

	// Specify the path for the new file
	newGenesisPath := filepath.Join(path, "genesis.json")
	newValidatorPath := filepath.Join(path, "priv_validator_key.json")

	_, err = os.Stat(newGenesisPath)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", newGenesisPath)
	} else {
		// Create a new file
		newGenesisFile, err := os.Create(newGenesisPath)
		if err != nil {
			log.Fatalf("failed to create new file: %v", err)
			panic(err)
		}
		defer newGenesisFile.Close()

		// Write the data to the new file
		_, err = newGenesisFile.Write(genesisData)
		if err != nil {
			log.Fatalf("failed to write data to new file: %v", err)
			panic(err)
		}
		log.Infof("New Cometbft Genesis file created successfully: %s\n", newGenesisPath)

	}

	_, err = os.Stat(newValidatorPath)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", newValidatorPath)
	} else {
		newValidatorFile, err := os.Create(newValidatorPath)
		if err != nil {
			log.Fatalf("failed to create new file: %v", err)
			panic(err)
		}
		defer newValidatorFile.Close()

		_, err = newValidatorFile.Write(validatorData)
		if err != nil {
			log.Fatalf("failed to write data to new file: %v", err)
			panic(err)
		}
		log.Infof("New Cometbft Validator file created successfully: %s\n", newValidatorPath)

	}
}

// InitCometbft initializes CometBFT for running a local sequencer.
func InitCometbft(defaultDir string, dataDirName string, binDirName string, configDirName string) {
	log.Info("Initializing CometBFT for running local sequencer:")
	cometbftDataPath := filepath.Join(defaultDir, dataDirName, ".cometbft")

	// verify that cometbft was downloaded and extracted to the correct location
	cometbftCmdPath := filepath.Join(defaultDir, binDirName, "cometbft")
	if !util.PathExists(cometbftCmdPath) {
		log.Error("Error: cometbft binary not found here", cometbftCmdPath)
		log.Error("\tCannot continue with initialization.")
		return
	}

	initCmdArgs := []string{"init", "--home", cometbftDataPath}
	initCmd := exec.Command(cometbftCmdPath, initCmdArgs...)

	log.Info("Running:", initCmd)

	_, err := initCmd.CombinedOutput()
	if err != nil {
		log.Error("Error executing command", initCmd, ":", err)
		return
	} else {
		log.Info("\tSuccess")
	}

	// copy the initialized genesis.json to the .cometbft directory
	initGenesisJsonPath := filepath.Join(defaultDir, configDirName, DefaultCometbftGenesisFilename)
	endGenesisJsonPath := filepath.Join(defaultDir, dataDirName, ".cometbft", "config", DefaultCometbftGenesisFilename)
	err = util.CopyFile(initGenesisJsonPath, endGenesisJsonPath)
	if err != nil {
		log.WithError(err).Error("Error copying CometBFT genesis.json file")
		return
	}
	log.Info("Copied genesis.json to", endGenesisJsonPath)

	// copy the initialized priv_validator_key.json to the .cometbft directory
	initPrivValidatorJsonPath := filepath.Join(defaultDir, configDirName, DefaultCometbftValidatorFilename)
	endPrivValidatorJsonPath := filepath.Join(defaultDir, dataDirName, ".cometbft", "config", DefaultCometbftValidatorFilename)
	err = util.CopyFile(initPrivValidatorJsonPath, endPrivValidatorJsonPath)
	if err != nil {
		log.WithError(err).Error("Error copying CometBFT priv_validator_key.json file")
		return
	}
	log.Info("Copied priv_validator_key.json to", endPrivValidatorJsonPath)

	// update the cometbft config.toml file to have the proper block time
	cometbftConfigPath := filepath.Join(defaultDir, dataDirName, ".cometbft/config/config.toml")
	oldValue := `timeout_commit = "1s"`
	newValue := `timeout_commit = "2s"`

	if err := replaceInFile(cometbftConfigPath, oldValue, newValue); err != nil {
		log.Error("Error updating the file:", cometbftConfigPath, ":", err)
		return
	} else {
		log.Info("Successfully updated", cometbftConfigPath)
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
		err := os.Rename(backupFilename, filename)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to rename temporary file to original: %w", err)
	}

	return nil
}
