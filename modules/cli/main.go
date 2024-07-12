package main

import (
	"github.com/astriaorg/astria-cli-go/modules/cli/cmd"
	// NOTE - must import the commands to register them
	_ "github.com/astriaorg/astria-cli-go/modules/cli/cmd/bundler"
	_ "github.com/astriaorg/astria-cli-go/modules/cli/cmd/devrunner"
	_ "github.com/astriaorg/astria-cli-go/modules/cli/cmd/sequencer"
	_ "github.com/astriaorg/astria-cli-go/modules/cli/cmd/sequencer/sudo"
)

func main() {
	cmd.Execute()
}
