# Astria Go

The `astria-go` cli is a tool designed to make local rollup development as
simple and dependency free as possible. It provides functionality to easily run
the Astria stack and interact with the Sequencer.

* [Available Commands](#available-commands)
* [Installation](#installation)
  * [Install and Run CLI from GitHub release](#install-and-run-cli-from-github-release)
  * [Locally Build and Run the CLI](#locally-build-and-run-the-cli)
* [Instances](#instances)
* [Development](#development)
  * [Testing](#testing)

## Available Commands

| Command                   | Description                                                                         |
|---------------------------|-------------------------------------------------------------------------------------|
| `version`                 | Prints the cli version.                                                             |
| `help`                    | Show help.                                                                          |
| `dev`                     | Root command for cli development functionality.                                     |
| `dev init`                | Downloads binaries and initializes the local environment.                           |
| `dev run`                 | Runs a minimal, local Astria stack.                                                 |
| `dev clean`               | Deletes the local data for the Astria stack.                                        |
| `dev clean all`           | Deletes the local data, downloaded binaries, and config files for the Astria stack. |
| `sequencer balances`      | Get the balances of an account on the Sequencer.                                    |
| `sequencer blockheight`   | Get the current block height of the Sequencer.                                      |
| `sequencer createaccount` | Generate an account for the Sequencer.                                              |
| `sequencer nonce`         | Get the current nonce for an account.                                               |
| `sequencer transfer`      | Get the current block height of the Sequencer.                                      |

## Installation

### Install and Run CLI from GitHub release

The CLI binaries are available for download from the
[releases page](https://github.com/astriaorg/astria-cli-go/releases). There are
binaries available for macOS and Linux, for both x86_64 and arm64 architectures.

```bash
# download the binary for your platform, e.g. macOS silicon
curl -L https://github.com/astriaorg/astria-cli-go/releases/download/v0.3.0/astria-cli-v0.3.0-darwin-arm64.tar.gz \
  --output astria-go.tar.gz
# extract the binary
tar -xzvf astria-go.tar.gz
# run the binary and check version
./astria-go version

# you can move the binary to a location in your PATH if you'd like
mv astria-go /usr/local/bin/
```

### Locally Build and Run the CLI

Dependencies: (only required for development)

* [GO](https://go.dev/doc/install)
* [just](https://github.com/casey/just)

```bash
git clone git@github.com:astriaorg/astria-cli-go.git
cd astria-cli-go
just build
just run "dev init"
just run "dev run"
```

This will download, configure, and run the following binaries of these
applications:

| App              | Version |
| ---------------- |---------|
| Cometbft         | v0.37.4 |
| Astria-Sequencer | v0.10.1 |
| Astria-Conductor | v0.13.1 |
| Astria-Composer  | v0.5.0  |

And place them in a `local-dev-astria` directory, along with several other
configuration files for everything.

The cli runs the minimum viable components for testing a rollup against the
Astria stack, allowing developers to confirm that their rollup interacts with
Astria's apis correctly.

## Instances

The `dev init`, `dev run`, and `dev clean` commands all have an optional
`--instance` flag. The value of this flag will be used as the directory name
where the rollup data will be stored. Now you can run many rollups while keeping
their configs and state data separate. If no value is provided, `default` is
used, i.e. `~/.astria/default`.

For example, if you run:

```bash
astria-go dev init
astria-go dev init --instance hello
astria-go dev init --instance world
```

You will see the following in the `~/.astria` directory:

```bash
.astria/
    default/
    hello/
    world/
```

Each of these directories will contain configs and binaries for
running the Astria stack. You can then update the `.env` files in the
`~/.astria/<instance name>/config-local/` or `~/.astria/<instance
name>/config-remote/` directories to suit your needs.

## Development

Requires go version 1.20 or newer.

You may also need to update your `gopls` settings in your editor for build tags
to allow for correct parsing of the build tags in the code. This will depend on
your IDE, but for VS Code you can open your settings and add:

```json
{
  "gopls": {
    "buildFlags": ["-tags=darwin arm64 amd64 linux"]
  }
}
```

The cli is built using [Cobra](https://github.com/spf13/cobra). Once you've
pulled the repo you can install the `cobra-cli` as follows to add new commands
for development:

```bash
# install cobra-cli
go install github.com/spf13/cobra-cli@latest
# add new command, e.g. `transfer`
cobra-cli add transfer
```

### Testing

```bash
# unit tests. some tests require tty.
go test ./...

# unit tests, skipping tests that require tty.
# this is useful for CI/CD pipelines.
go test ./... -skip TestProcessPane

# integration tests. requires running geth + cometbft + astria stack
# build binary used for testing
go build -o ./bin/astria-go-testy
# run integration tests. requires -tag
go test ./integration_tests -tags=integration_tests
# cleanup binary
rm ./bin/astria-go-testy

# of course, there are just commands to make this easier
just test
just t

just test-integration
just ti

just test-all
just ta
```
