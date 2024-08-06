# Astria Go Monorepo

This repository contains all the Go packages that are used within the Astria
platform.

* [bech32m](./modules/bech32m/README.md) - a utility package that is a simple
  facade for
  [btcutil/bech32m](https://github.com/btcsuite/btcd/tree/master/btcutil/bech32)
* [cli](./modules/cli/README.md) - The Astria CLI is a command line
  interface for the Astria platform. It is used to interact with the Astria
  platform and manage your rollup projects. It has two main commands:
  * `dev` - This command is used to start the development server for your
    rollup project.
  * `sequencer` - This command has subcommands to interact with a local or
    remote Sequencer network, e.g. `createaccount`, `transfer`, etc.
* [go-sequencer-client](./modules/go-sequencer-client/README.md) - a client library for
  interacting with a Sequencer network. The Astria CLI uses this package to
  interact with the Sequencer. It is also used to develop rollups on top of
  Astria
