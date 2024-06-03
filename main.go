package main

import (
	"github.com/astria/astria-cli-go/cmd"
	// NOTE - must import the commands to register them
	_ "github.com/astria/astria-cli-go/cmd/devrunner"
	_ "github.com/astria/astria-cli-go/cmd/sequencer"
	_ "github.com/astria/astria-cli-go/cmd/sequencer/sudo"
)

func main() {
	cmd.Execute()
}
