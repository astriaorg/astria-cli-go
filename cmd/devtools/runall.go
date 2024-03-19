package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astria/astria-cli-go/internal/processrunner"
	"github.com/astria/astria-cli-go/internal/ui"
	"github.com/spf13/cobra"
)

// runallCmd represents the runall command
var runallCmd = &cobra.Command{
	Use:   "runall",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runall()
	},
}

func init() {
	rootCmd.AddCommand(runallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runall() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
		return
	}
	defaultDir := filepath.Join(homePath, ".astria")

	// Load the .env file and get the environment variables
	envPath := filepath.Join(defaultDir, "local-dev-astria/.env")
	environment := loadAndGetEnvVariables(envPath)
	//fmt.Printf("Environment: %v\n", environment)

	sequencerBinPath := filepath.Join(homePath, ".astria/local-dev-astria/astria-sequencer")
	seqProcRunner := processrunner.NewProcessRunner("Sequencer", sequencerBinPath, environment, nil)

	// shouldStart acts as a control channel to start this first process
	shouldStart := make(chan bool)
	seqProcRunner, err = seqProcRunner.Start(shouldStart)
	shouldStart <- true
	if err != nil {
		fmt.Println("Error running sequencer:", err)
		panic(err)
	}

	cometBinPath := filepath.Join(homePath, ".astria/local-dev-astria/cometbft")
	cometDataPath := filepath.Join(homePath, ".astria/data/.cometbft")
	cometArgs := []string{"node", "--home", cometDataPath}
	cometProcRunner := processrunner.NewProcessRunner("Comet BFT", cometBinPath, environment, cometArgs)
	cometProcRunner, err = cometProcRunner.Start(seqProcRunner.IsRunning)
	if err != nil {
		fmt.Println("Error running composer:", err)
		panic(err)
	}

	composeBinPath := filepath.Join(homePath, ".astria/local-dev-astria/astria-composer")
	compProcRunner := processrunner.NewProcessRunner("Composer", composeBinPath, environment, nil)
	compProcRunner, err = compProcRunner.Start(cometProcRunner.IsRunning)
	if err != nil {
		fmt.Println("Error running composer:", err)
		panic(err)
	}

	conductorBinPath := filepath.Join(homePath, ".astria/local-dev-astria/astria-conductor")
	condProcRunner := processrunner.NewProcessRunner("Conductor", conductorBinPath, environment, nil)
	condProcRunner, err = condProcRunner.Start(compProcRunner.IsRunning)
	if err != nil {
		fmt.Println("Error running conductor:", err)
		panic(err)
	}

	runners := []*processrunner.ProcessRunner{seqProcRunner, cometProcRunner, compProcRunner, condProcRunner}
	app := ui.NewApp(runners)

	app.Start()
}
