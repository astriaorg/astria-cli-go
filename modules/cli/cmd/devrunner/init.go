package devrunner

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/config"

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

	flagHandler := cmd.CreateCliFlagHandler(initCmd, cmd.EnvPrefix)
	flagHandler.BindStringFlag("local-network-name", config.DefaultLocalNetworkName, "Set the local network name for the instance. This is used to set the chain ID in the CometBFT genesis.json file.")
	flagHandler.BindStringFlag("local-native-denom", config.DefaultLocalNativeDenom, "Set the default denom for the local instance. This is used to set the 'native_asset_base_denomination' and 'allowed_fee_assets' in the CometBFT genesis.json file.")
	flagHandler.BindStringFlag("rollup-name", config.DefaultRollupName, "Set the default rollup name for the local instance. This is used to set the 'astria_composer_rollups' in the base-config.toml file.")
}

func runInitialization(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	localNetworkName := flagHandler.GetValue("local-network-name")
	config.IsSequencerChainIdValidOrPanic(localNetworkName)

	rollupName := flagHandler.GetValue("rollup-name")
	config.IsSequencerChainIdValidOrPanic(rollupName)

	localDenom := flagHandler.GetValue("local-native-denom")

	// TODO: make the default home dir configurable
	homeDir := cmd.GetUserHomeDirOrPanic()
	instanceDir := filepath.Join(homeDir, ".astria", instance)

	// paths must be absolute
	logsDir := filepath.Join(homeDir, ".astria", instance, config.LogsDirName)
	localBinDir := filepath.Join(homeDir, ".astria", instance, config.BinariesDirName)
	networksConfigPath := filepath.Join(homeDir, ".astria", instance, config.DefaultNetworksConfigName)
	tuiConfigPath := filepath.Join(homeDir, ".astria", config.DefaultTUIConfigName)
	configDir := filepath.Join(homeDir, ".astria", instance, config.DefaultConfigDirName)
	baseConfigPath := filepath.Join(homeDir, ".astria", instance, config.DefaultConfigDirName, config.DefaultBaseConfigName)

	log.Info("Creating new instance in:", instanceDir)
	cmd.CreateDirOrPanic(instanceDir)
	cmd.CreateDirOrPanic(configDir)

	log.Info("Binary files for locally running a services placed in: ", localBinDir)
	cmd.CreateDirOrPanic(localBinDir)
	cmd.CreateDirOrPanic(logsDir)

	genericBinariesDir := filepath.Join("~", ".astria", instance, config.BinariesDirName)
	config.CreateNetworksConfig(networksConfigPath, genericBinariesDir, localNetworkName, rollupName, localDenom)
	networkConfigs := config.LoadNetworkConfigsOrPanic(networksConfigPath)

	config.CreateTUIConfig(tuiConfigPath)

	config.CreateBaseConfig(baseConfigPath, instance, localNetworkName, rollupName, localDenom)

	config.CreateComposerDevPrivKeyFile(configDir)

	config.RecreateCometbftAndSequencerGenesisData(configDir, localNetworkName, localDenom)

	// download and unpack all services for all networks
	for label := range networkConfigs.Configs {
		purpleANSI := "\033[35m"
		resetANSI := "\033[0m"
		log.Info(fmt.Sprint("--Downloading binaries for network: ", purpleANSI, label, resetANSI))
		for _, bin := range networkConfigs.Configs[label].Services {
			downloadAndUnpack(bin.DownloadURL, bin.Version, bin.Name, localBinDir)
		}
	}

	// create the data directory for cometbft and sequencer
	dataPath := filepath.Join(homeDir, ".astria", instance, config.DataDirName)
	cmd.CreateDirOrPanic(dataPath)
	config.InitCometbft(instanceDir, config.DataDirName, config.BinariesDirName, config.MainnetCometbftVersion, config.DefaultConfigDirName)

	log.Infof("Initialization of instance \"%s\" completed successfuly.", instance)

}

// downloadFile downloads a file from the specified URL to the given local path.
func downloadFile(url, filepath string) error {
	// get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// write downloaded data to a file
	_, err = io.Copy(out, resp.Body)
	return err
}

// extractTarGz extracts a .tar.gz file to dest.
func extractTarGz(dest string, version string, gzipStream io.Reader) error {
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
		target := filepath.Join(dest, header.Name+"-"+version)

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

// downloadAndUnpack downloads a file from the specified URL, extracts it, and
// places it in the given path.
//
// Panics if the download, extraction, deletion of the .tar.gz, or placement of
// the file fails.
func downloadAndUnpack(url, version, packageName, placePath string) {
	if url == "" {
		log.Infof("No source URL provided for %s. Skipping download.\n", packageName)
		return
	}

	// check if the file already exists
	if _, err := os.Stat(filepath.Join(placePath, packageName+"-"+version)); err == nil {
		log.Infof("%s already exists. Skipping download.\n", packageName)
		return
	}
	log.Infof("Downloading: %s, %s, %s\n", packageName, version, url)

	// download the file
	dest := filepath.Join(placePath, packageName+"-"+version+".tar.gz")
	if err := downloadFile(url, dest); err != nil {
		panic(err)
	}
	// open the downloaded .tar.gz file
	file, err := os.Open(dest)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// extract the contents
	if err := extractTarGz(placePath, version, file); err != nil {
		panic(err)
	}

	// delete the .tar.gz file
	// TODO: should this be configurable?
	err = os.Remove(dest)
	if err != nil {
		log.Fatalf("Failed to delete downloaded %s.tar.gz file: %v", packageName, err)
		panic(err)
	}
	log.Infof("%s downloaded and extracted successfully.\n", packageName)
}
