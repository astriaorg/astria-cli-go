package devtools

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/astria/astria-cli-go/cmd"
	"github.com/astria/astria-cli-go/cmd/devtools/config"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:    "init",
	Short:  "Initializes the local development environment.",
	Long:   `The init command will download the necessary binaries, create new directories for file organisation, and create an environment file for running a minimal Astria stack locally.`,
	PreRun: cmd.SetLogLevel,
	Run:    runInitialization,
}

func init() {
	devCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("instance", "i", config.DefaultInstanceName, "Used to set the directory name in ~/.astria to enable running separate instances of the sequencer stack.")
	initCmd.Flags().String("local-network-name", "sequencer-test-chain-0", "Set the network name for the local instance. This is used to set the chain ID in the CometBFT genesis.json file.")
	initCmd.Flags().String("local-default-denom", "nria", "Set the default denom for the local instance. This is used to set the 'native_asset_base_denomination' and 'allowed_fee_assets' in the CometBFT genesis.json file.")
}

func runInitialization(c *cobra.Command, args []string) {
	// Get the instance name from the -i flag or use the default
	instance := c.Flag("instance").Value.String()
	config.IsInstanceNameValidOrPanic(instance)

	localNetworkName := c.Flag("local-network-name").Value.String()
	config.IsSequencerChainIdValidOrPanic(localNetworkName)

	localDefaultDenom := c.Flag("local-default-denom").Value.String()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error("error getting home dir:", err)
		return
	}
	// TODO: make the default home dir configurable
	defaultDir := filepath.Join(homeDir, ".astria")
	instanceDir := filepath.Join(defaultDir, instance)

	log.Info("Creating new instance in:", instanceDir)
	cmd.CreateDirOrPanic(instanceDir)

	networksConfigPath := filepath.Join(defaultDir, instance, config.DefualtNetworksConfigName)
	config.CreateNetworksConfig(networksConfigPath, localNetworkName, localDefaultDenom)

	configDirPath := filepath.Join(instanceDir, config.DefaultConfigDirName)
	cmd.CreateDirOrPanic(configDirPath)

	baseConfigPath := filepath.Join(configDirPath, config.DefualtBaseConfigName)
	config.CreateBaseConfig(baseConfigPath, instance)

	// create the local config and env files
	// configPath := filepath.Join(instanceDir, config.ConfigDirName)

	// cmd.CreateDirOrPanic(configPath)
	config.RecreateCometbftAndSequencerGenesisData(configDirPath, localNetworkName, localDefaultDenom)

	// create the local bin directory for downloaded binaries
	localBinPath := filepath.Join(instanceDir, config.BinariesDirName)
	log.Info("Binary files for locally running a sequencer placed in: ", localBinPath)
	cmd.CreateDirOrPanic(localBinPath)
	for _, bin := range config.Binaries {
		downloadAndUnpack(bin.Url, bin.Name, localBinPath)
	}

	// create the data directory for cometbft and sequencer
	dataPath := filepath.Join(instanceDir, config.DataDirName)
	cmd.CreateDirOrPanic(dataPath)
	config.InitCometbft(instanceDir, config.DataDirName, config.BinariesDirName, config.DefaultConfigDirName)

	log.Infof("Initialization of instance \"%s\" completed successfuly.", instance)

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
				cmd.CreateDirOrPanic(target)
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

func downloadAndUnpack(url string, packageName string, placePath string) {
	// Check if the file already exists
	if _, err := os.Stat(filepath.Join(placePath, packageName)); err == nil {
		log.Infof("%s already exists. Skipping download.\n", packageName)
		return
	}
	log.Infof("Downloading: (%s, %s)\n", packageName, url)

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
		panic(err)
	}
	log.Infof("%s downloaded and extracted successfully.\n", packageName)
}
