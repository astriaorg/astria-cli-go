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
	flagHandler.BindStringFlag("local-native-denom", config.LocalNativeDenom, "Set the default denom for the local instance. This is used to set the 'native_asset_base_denomination' and 'allowed_fee_assets' in the CometBFT genesis.json file.")
}

func runInitialization(c *cobra.Command, args []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	localNetworkName := flagHandler.GetValue("local-network-name")
	config.IsSequencerChainIdValidOrPanic(localNetworkName)

	localDefaultDenom := flagHandler.GetValue("local-native-denom")

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

	// create the local bin directory for downloaded binaries
	localBinPath := filepath.Join(instanceDir, config.BinariesDirName)
	log.Info("Binary files for locally running a sequencer placed in: ", localBinPath)
	cmd.CreateDirOrPanic(localBinPath)

	networksConfigPath := filepath.Join(defaultDir, instance, config.DefaultNetworksConfigName)
	config.CreateNetworksConfig(localBinPath, networksConfigPath, localNetworkName, localDefaultDenom)
	networkConfigs := config.LoadNetworkConfigsOrPanic(networksConfigPath)

	configDirPath := filepath.Join(instanceDir, config.DefaultConfigDirName)
	cmd.CreateDirOrPanic(configDirPath)

	baseConfigPath := filepath.Join(configDirPath, config.DefaultBaseConfigName)
	config.CreateBaseConfig(baseConfigPath, instance)

	config.CreateComposerDevPrivKeyFile(configDirPath)

	config.RecreateCometbftAndSequencerGenesisData(configDirPath, localNetworkName, localDefaultDenom)

	// download and unpack all services for all networks
	for label := range networkConfigs.Configs {
		purpleANSI := "\033[35m"
		resetANSI := "\033[0m"
		log.Info(fmt.Sprint("--Downloading binaries for network: ", purpleANSI, label, resetANSI))
		for _, bin := range networkConfigs.Configs[label].Services {
			downloadAndUnpack(bin.DownloadURL, bin.Version, bin.Name, localBinPath)
		}
	}

	// create the data directory for cometbft and sequencer
	dataPath := filepath.Join(instanceDir, config.DataDirName)
	cmd.CreateDirOrPanic(dataPath)
	config.InitCometbft(instanceDir, config.DataDirName, config.BinariesDirName, config.CometbftVersion, config.DefaultConfigDirName)

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
