package devrunner

import (
	"fmt"
	"path/filepath"

	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner/config"
	"github.com/spf13/cobra"
)

// VersionCmd represents the sequencer command
var VersionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the version of the services used by the CLI.",
	Aliases: []string{"versions"},
	Run:     seqVersionCmdHandler,
}

func seqVersionCmdHandler(c *cobra.Command, _ []string) {
	flagHandler := cmd.CreateCliFlagHandler(c, cmd.EnvPrefix)

	homeDir := cmd.GetUserHomeDirOrPanic()
	astriaDir := filepath.Join(homeDir, ".astria")
	instance := flagHandler.GetValue("instance")
	config.IsInstanceNameValidOrPanic(instance)

	network := flagHandler.GetValue("network")

	networksConfigPath := filepath.Join(astriaDir, instance, config.DefaultNetworksConfigName)
	networkConfigs := config.LoadNetworkConfigsOrPanic(networksConfigPath)

	networkOverrides := networkConfigs.Configs[network]

	fmt.Println("Preset Service Versions:")
	// get longest service name for padding
	longestServiceName := 0
	for _, service := range networkOverrides.Services {
		if len(service.Name) > longestServiceName {
			longestServiceName = len(service.Name)
		}
	}
	// print service versions
	for _, service := range networkOverrides.Services {
		fmt.Printf("%-*s: %s\n", longestServiceName, service.Name, service.Version)
	}
}

func init() {
	devCmd.AddCommand(VersionCmd)

	flagHandler := cmd.CreateCliFlagHandler(VersionCmd, cmd.EnvPrefix)
	flagHandler.BindStringFlag("network", config.DefaultTargetNetwork, "Select the network to print the service versions for.")
}
