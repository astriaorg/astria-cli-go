package main

import (
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	// NOTE - must import the commands to register them
	_ "github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner"
	_ "github.com/astriaorg/astria-cli-go/modules/cli/cmd/sequencer"
)

func main() {
	cmd.Execute()
}
