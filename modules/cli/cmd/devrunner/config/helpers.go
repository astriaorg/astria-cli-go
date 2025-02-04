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
	"sort"
	"strings"
	"unicode"

	util "github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/utilities"

	log "github.com/sirupsen/logrus"
)

// IsInstanceNameValidOrPanic checks if the instance name is valid.
//
// Panics if the instance name is not valid.
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

// IsSequencerChainIdValidOrPanic checks if the instance name is valid.
//
// Panics if the instance name is not valid.
func IsSequencerChainIdValidOrPanic(id string) {
	if len(id) < 1 || len(id) > 50 {
		log.Errorf("Invalid sequencer chain id length: %s", id)
		err := fmt.Errorf("invalid sequencer chain id: '%s'. The ChainId length must be within the range [1,50]", id)
		panic(err)
	}

	re, err := regexp.Compile(`^[a-zA-Z0-9\-_\.]+$`)
	if err != nil {
		log.WithError(err).Error("Error compiling regex")
		panic(err)
	}
	if !re.MatchString(id) {
		log.Errorf("Invalid sequencer chain id: %s", id)
		err := fmt.Errorf("invalid sequencer chain id: '%s'. The ChainId can only contain lowercase and uppercase letters, numerical digits, and the characters '-', '_', and '.'", id)
		panic(err)
	}
}

//go:embed composer_dev_priv_key
var embeddedDevPrivKey embed.FS

// CreateComposerDevPrivKeyFile creates a new composer_dev_priv_key file in the specified directory.
func CreateComposerDevPrivKeyFile(dir string) {
	dir = util.ShellExpand(dir)
	// read the content from the embedded file
	devPrivKeyData, err := fs.ReadFile(embeddedDevPrivKey, "composer_dev_priv_key")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
		panic(err)
	}

	// specify the path for the new file
	newDevPrivKeyPath := filepath.Join(dir, "composer_dev_priv_key")

	_, err = os.Stat(newDevPrivKeyPath)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", newDevPrivKeyPath)
	} else {
		// create a new file
		newDevPrivKeyFile, err := os.Create(newDevPrivKeyPath)
		if err != nil {
			log.Fatalf("failed to create new file: %v", err)
			panic(err)
		}
		defer newDevPrivKeyFile.Close()

		// write the data to the new file
		_, err = newDevPrivKeyFile.Write(devPrivKeyData)
		if err != nil {
			log.Fatalf("failed to write data to new file: %v", err)
			panic(err)
		}
		log.Infof("New composer_dev_priv_key file created successfully: %s\n", newDevPrivKeyPath)
	}
}

//go:embed genesis.json
var embeddedCometbftGenesisFile embed.FS

//go:embed priv_validator_key.json
var embeddedCometbftValidatorFile embed.FS

// RecreateCometbftAndSequencerGenesisData creates a new CometBFT genesis.json
// and priv_validator_key.json file at the specified path.
//   - path: the path to the directory where the new files will be created.
//   - localNetworkName: the name of the local sequencer network.
//   - localNativeDenom: the native denomination for the local sequencer network.
//
// Panics if the files cannot be created.
func RecreateCometbftAndSequencerGenesisData(path, localNetworkName, localNativeDenom string) {
	path = util.ShellExpand(path)
	// read the content from the embedded file
	genesisData, err := fs.ReadFile(embeddedCometbftGenesisFile, "genesis.json")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
		panic(err)
	}
	// unmarshal JSON into a map to update sequencer chain id
	var data map[string]interface{}
	if err := json.Unmarshal(genesisData, &data); err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}
	// update chain id and default denom and convert back to bytes
	data["chain_id"] = localNetworkName
	if appState, ok := data["app_state"].(map[string]interface{}); ok {
		appState["native_asset_base_denomination"] = localNativeDenom
		appState["allowed_fee_assets"] = []interface{}{localNativeDenom}
		data["app_state"] = appState
	} else {
		log.Println("Error: Expected map[string]interface{} for 'app_state'")
	}
	genesisData, err = json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling updated data to JSON: %s", err)
	}

	// read the content from the embedded file
	validatorData, err := fs.ReadFile(embeddedCometbftValidatorFile, "priv_validator_key.json")
	if err != nil {
		log.Fatalf("failed to read embedded file: %v", err)
		panic(err)
	}

	// specify the path for the new file
	newGenesisPath := filepath.Join(path, "genesis.json")
	newValidatorPath := filepath.Join(path, "priv_validator_key.json")

	_, err = os.Stat(newGenesisPath)
	if err == nil {
		log.Infof("%s already exists. Skipping initialization.\n", newGenesisPath)
	} else {
		// create a new file
		newGenesisFile, err := os.Create(newGenesisPath)
		if err != nil {
			log.Fatalf("failed to create new file: %v", err)
			panic(err)
		}
		defer newGenesisFile.Close()

		// write the data to the new file
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
func InitCometbft(defaultDir string, dataDirName string, binDirName string, binVersion string, configDirName string) {
	defaultDir = util.ShellExpand(defaultDir)
	log.Info("Initializing CometBFT for running local sequencer:")
	cometbftDataPath := filepath.Join(defaultDir, dataDirName, ".cometbft")

	// verify that cometbft was downloaded and extracted to the correct location
	cometbftCmdPath := filepath.Join(defaultDir, binDirName, "cometbft-v"+binVersion)
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

	if err := ReplaceInFile(cometbftConfigPath, oldValue, newValue); err != nil {
		log.Error("Error updating the file:", cometbftConfigPath, ":", err)
		return
	} else {
		log.Info("Successfully updated", cometbftConfigPath)
	}
}

// ReplaceInFile replaces oldValue with newValue in the file at filename.
// it is used here to update the block time in the cometbft config.toml file.
func ReplaceInFile(filename, oldValue, newValue string) error {
	// read the original file.
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read the file: %w", err)
	}

	// perform the replacement.
	modifiedContent := strings.ReplaceAll(string(content), oldValue, newValue)

	// write the modified content to a new temporary file.
	tmpFilename := filename + ".tmp"
	if err := os.WriteFile(tmpFilename, []byte(modifiedContent), 0666); err != nil {
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}

	// rename the original file to filename.bak.
	backupFilename := filename + ".bak"
	if err := os.Rename(filename, backupFilename); err != nil {
		return fmt.Errorf("failed to rename original file to backup: %w", err)
	}

	// rename the temporary file to the original file name.
	if err := os.Rename(tmpFilename, filename); err != nil {
		// attempt to restore the original file if renaming fails.
		err := os.Rename(backupFilename, filename)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to rename temporary file to original: %w", err)
	}

	// remove the backup file.
	backupFile, err := os.Open(backupFilename)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	if err := backupFile.Close(); err != nil {
		return fmt.Errorf("failed to close backup file: %w", err)
	}
	if err := os.Remove(backupFilename); err != nil {
		return fmt.Errorf("failed to remove backup file: %w", err)
	}

	return nil
}

// MergeConfigs merges two or more slices of "key=value" strings into a single slice.
// The slices are merged in order, with later slices overwriting earlier ones.
func MergeConfigs(configs ...[]string) []string {
	mergedMap := make(map[string]string)

	// helper function to add slices to the map
	addSliceToMap := func(slice []string) {
		for _, item := range slice {
			keyVal := strings.SplitN(item, "=", 2)
			if len(keyVal) != 2 {
				continue // skip any items that don't correctly split into two parts
			}
			key, value := keyVal[0], keyVal[1]
			mergedMap[key] = value
		}
	}

	for _, config := range configs {
		addSliceToMap(config)
	}

	// convert the map back to a slice
	var result []string
	for key, value := range mergedMap {
		result = append(result, key+"="+value)
	}
	sort.Strings(result)

	return result
}

// LogEnv logs the configuration to the cli log file.
func LogEnv(env []string) {
	log.Debug("Environment:")
	for _, item := range env {
		log.Debug(item)
	}
}

// validateServiceLogLevelOrPanic validates the service log level and panics if
// it is invalid. The valid log levels are: debug, info, error.
func validateServiceLogLevelOrPanic(logLevel string) {
	switch logLevel {
	case "debug", "info", "error":
		return
	default:
		log.WithField("service-log-level", logLevel).Fatal("Invalid service log level. Must be one of: 'debug', 'info', 'error'")
		panic("Invalid service log level")
	}

}

// GetServiceLogLevelOverrides returns a slice of strings that can be used to
// update the log level for the Astria services.
//
// The env var log levels that are returned are:
//   - ASTRIA_SEQUENCER_LOG
//   - ASTRIA_COMPOSER_LOG
//   - ASTRIA_CONDUCTOR_LOG
//
// Panics if the service log level is not one of the following: debug, info, error.
func GetServiceLogLevelOverrides(serviceLogLevel string) []string {
	validateServiceLogLevelOrPanic(serviceLogLevel)
	serviceLogLevelOverrides := []string{
		"ASTRIA_SEQUENCER_LOG=\"astria_sequencer=" + serviceLogLevel + "\"",
		"ASTRIA_COMPOSER_LOG=\"astria_composer=" + serviceLogLevel + "\"",
		"ASTRIA_CONDUCTOR_LOG=\"astria_conductor=" + serviceLogLevel + "\"",
	}
	return serviceLogLevelOverrides
}

// IsValidDenom checks if the input string is a valid denomination.
//
// A valid denomination is a string that contains only letters.
//
// Panics if the input string is not a valid denomination.
func IsValidDenomOrPanic(denom string) {
	denom = strings.ToLower(denom)

	for _, r := range denom {
		if !unicode.IsLetter(r) {
			log.Error("Error validating denomination:", denom, "Denominations must contain only letters.")
			panic("Invalid denomination: " + denom + ", denominations must contain only letters.")
		}
	}
}
