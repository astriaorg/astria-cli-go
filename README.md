# Astria Go

The `astria-go` cli is a tool designed to make local rollup development as
simple and dependency free as possible. It provides functionality to easily run
the Astria stack and interact with the Sequencer.

* [Available Commands](#available-commands)
* [Installation](#installation)
  * [Install from GitHub release](#install-from-github-release)
  * [Build Locally from Source](#build-locally-from-source)
* [Running the Astria Sequencer](#running-the-astria-sequencer)
  * [Initialize Configuration](#initialize-configuration)
  * [Usage](#usage)
    * [Run a Local Sequencer](#run-a-local-sequencer)
    * [Run Against a Remote Sequencer](#run-against-a-remote-sequencer)
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

### Install from GitHub release

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

### Build Locally from Source

Dependencies: (only required for development)

* [GO](https://go.dev/doc/install)
* [just](https://github.com/casey/just)

```bash
# checkout repo
git clone git@github.com:astriaorg/astria-cli-go.git
cd astria-cli-go

# run build command
just build

# check the version
just run "version"
# or
go run main.go version
```

## Running the Astria Sequencer

### Initialize Configuration

```bash
astria-go dev init
```

The `init` command downloads binaries and generates environment and
configuration files.

The following files are generated:

* TODO

The following binaries are downloaded:

| App              | Version |
|------------------|---------|
| Cometbft         | v0.37.4 |
| Astria-Sequencer | v0.10.1 |
| Astria-Conductor | v0.13.1 |
| Astria-Composer  | v0.5.0  |

### Usage

The cli runs the minimum viable components for testing a rollup against the
Astria stack, allowing developers to confirm that their rollup interacts with
Astria's APIs correctly.

You can choose to run the Sequencer locally, or you can run the stack against
the remote Sequencer. You may also run local binaries instead of downloaded
pre-built binaries.

#### Run a Local Sequencer

The simplest way to run Astria:

```bash
astria-go dev run
```

This will spin up a Sequencer (Cometbft and Astria-Sequencer), a Conductor,
and a Composer -- all on your local machine, using pre-built binaries of the
dependencies.

> NOTE: Running a local Sequencer is the default behavior of `dev run` command.
> Thus, `astria-go dev run` is effectively an alias of
> `astria-go dev run --local`.

#### Run Against a Remote Sequencer

If you want to run Composer and Conductor locally against a remote Astria
Sequencer:

```bash
astria-go dev run --remote
```

Using the `--remote` flag, the cli will handle configuration of the components
running on your local machine, but you will need to create an account on the
remote sequencer. More details can be
[found here](https://docs.astria.org/developer/tutorials/1-using-astria-go-cli#setup-and-run-the-local-astria-components-to-communicate-with-the-remote-sequencer).

### Run Custom Binaries

You can also run components of Astria from a local monorepo during development
of Astria core itself. For example if you are developing a new feature in the
[`astria-conductor` crate](https://github.com/astriaorg/astria/tree/main/crates/astria-conductor)
in the [Astria mono repo](https://github.com/astriaorg/astria) you can use the
cli to run your locally compiled Conductor with the other components using the
`--conductor-path` flag:

```bash
astria-go dev run --local \
  --conductor-path <absolute path to the Astria mono repo>/target/debug/astria-conductor
```

This will run Composer, Cometbft, and Sequencer using the downloaded pre-built
binaries, while using a locally built version of the Conductor binary. You can
swap out some or all binaries for the Astria stack with their appropriate flags:

```bash
astria-go dev run --local \
  --sequencer-path <sequencer bin path> \
  --cometbft-path <cometbft bin path> \
  --composer-path <composer bin path> \
  --conductor-path <conductor bin path>
```

If additional configuration is required, you can update the `.env` files in
`~/.astria/<instance>/config-local/` or
`~/.astria/<instance>/config-remote/` based on your needs.

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

Requires go version 1.22 or newer.

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
