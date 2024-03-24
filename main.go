package main

import (
	"github.com/astria/astria-cli-go/cmd"
	// NOTE - must import the commands to register them
	_ "github.com/astria/astria-cli-go/cmd/devtools"
	_ "github.com/astria/astria-cli-go/cmd/sequencer"
)

func main() {
	cmd.Execute()
}
