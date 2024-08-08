# Astria Go

The `astria-go` cli is a tool designed to make local rollup development as
simple and dependency free as possible. It provides functionality to easily run
the Astria stack and interact with the Sequencer.

* [Installation](#installation)
  * [Install from GitHub release](#install-from-github-release)
  * [Install from Nightly release](#install-nightly-release)
  * [Build Locally from Source](#build-locally-from-source)
* [Running the Astria Sequencer](#running-the-astria-sequencer)
  * [Initialize Configuration](#initialize-configuration)
  * [Usage](#usage)
    * [Run a Local Sequencer](#run-a-local-sequencer)
    * [Run Against a Remote Sequencer](#run-against-a-remote-sequencer)
  * [Run Custom Binaries](#run-custom-binaries)
  * [Interact with the Sequencer](#interact-with-the-sequencer)
* [Instances](#instances)
* [Configuration](#configuration)
  * [Set Service Env Vars with the
    `base-config.toml`](#set-service-env-vars-with-the-base-configtoml)
  * [Configure the Services and Sequencer Networks with
    `networks-config.toml`](#configure-the-services-and-sequencer-networks-with-networks-configtoml)
    * [Choosing Specific Services](#choosing-specific-services)
    * [Adding a New Network](#adding-a-new-network)
  * [Interact with a Sequencer Network using the `sequencer-networks-config.toml`](#interact-with-a-sequencer-network-using-the-sequencer-networks-configtoml)
* [Development](#development)
  * [TUI Logs](#tui-logs)
  * [Testing](#testing)

## Installation

### Install From GitHub Release

The latest release for the CLI can be found here
[v0.12.0](https://github.com/astriaorg/astria-cli-go/releases/tag/v0.12.0).
There are binaries available for macOS (arm64 and x86_64) and Linux(x86_64) architectures.

```bash
# download the binary for your platform, e.g. macOS silicon
curl -L https://github.com/astriaorg/astria-cli-go/releases/download/v0.12.0/astria-go-v0.12.0-darwin-arm64.tar.gz \
  --output astria-go.tar.gz
# extract the binary
tar -xzvf astria-go.tar.gz
# run the binary and check version
./astria-go version

# you can move the binary to a location in your PATH if you'd like
mv astria-go /usr/local/bin/
```

### Install Nightly Release

The nightly releases for the cli can be found on the cli [releases
page](https://github.com/astriaorg/astria-cli-go/releases).
Grab the

```bash
export NIGHTLY_URL="download url of the build you need"
curl -L $NIGHTLY_URL > astria-cli.tar.gz
tar -xvzf astria-cli.tar.gz
mv astria-go /usr/local/bin/
astria-go version

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
# will return "development" because you built locally
just run "version" 
# or
go run main.go version
```

## Running the Astria Sequencer

### Initialize Configuration

```bash
astria-go dev init
```

The `init` command downloads binaries, generates configuration files, downloads
the service binaries, and initializes CometBFT. These files will be organized
and place in the `~/.astria` directory within an `/<instance>` directory. See
the section on [Instances](#instances) for more details.

The `init` command will also run the initialization steps required by CometBFT.

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
astria-go dev run --network local
```

This will spin up a Sequencer (Cometbft and Astria-Sequencer), a Conductor,
and a Composer all on your local machine, using pre-built binaries of the
dependencies.

#### Run Against a Remote Sequencer

If you want to run Composer and Conductor locally against a remote Astria
Sequencer:

```bash
# Run against the Astria Dusk dev net 
astria-go dev run --network dusk
```

When using the `--network` flag to target a remote sequencer, the cli will
handle configuration of the components running on your local machine, but you
will need to create an account on the remote sequencer. More details can be
[found here](https://docs.astria.org/developer/tutorials/1-using-astria-go-cli#setup-and-run-the-local-astria-components-to-communicate-with-the-remote-sequencer).

### Run Custom Binaries

You can also run components of Astria from a local monorepo during development
of Astria core itself. For example if you are developing a new feature in the
[`astria-conductor` crate](https://github.com/astriaorg/astria/tree/main/crates/astria-conductor)
in the [Astria mono repo](https://github.com/astriaorg/astria) you can use the
cli to run your locally compiled Conductor with the other components using the
`--conductor-path` flag:

```bash
astria-go dev run --network local \
  --conductor-path <absolute path to the Astria mono repo>/target/debug/astria-conductor
```

Or update the local path to the service you want to change in the networks
config: `~/.astria/<instance>/networks-config.toml`

```toml
[networks.local.services.conductor]
name = 'astria-conductor'
version = 'dev' # update the version to be relevant
download_url = '' # don't need a download url for local binaries
local_path = 'path to your local conductor'
```

This will run Composer, Cometbft, and Sequencer using the downloaded pre-built
binaries, while using a locally built version of the Conductor binary. You can
swap out some or all binaries for the Astria stack with their appropriate flags:

```bash
astria-go dev run --network local \
  --sequencer-path <sequencer bin path> \
  --cometbft-path <cometbft bin path> \
  --composer-path <composer bin path> \
  --conductor-path <conductor bin path>
```

Or by updating the specific service in the `networks-config.toml`.

### Interact with the Sequencer

You can use the `astria-go sequencer [command]` commands to interact with the
sequencer. Run:

```bash
astria-go sequencer -h
```

To see a full list of commands available.

You can also use the `--network` flag on the relevant commands to automatically
configure which sequencer network the commands will be run against.

## Instances

The `dev` commands all have an optional `--instance` flag. The value of this
flag will be used as the directory name where the rollup data will be stored.
Now you can run many rollups while keeping their configs and state data
separate. If no value is provided, `default` is used, i.e. `~/.astria/default`.

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
running the Astria stack.

## Configuration

There are three different configuration files that the cli uses:

* `base-config.toml`
* `networks-config.toml`
* `sequencer-networks-config.toml`

The `base-config.toml` and `networks-config.toml` are used with the `astria-go
dev` commands, and the `sequencer-networks-config.toml` is used with the
`astria-go sequencer` commands.

### Set Service Env Vars with the `base-config.toml`

The base config can be found in the `~/.astria/<instance>/config` directory. It
contains all of the configuration settings for the services that will be run by
the CLI. The file is parsed from toml into environment variables that are then
passed to all services.
If you need to add or change any of the settings for the services you can add a
new setting to the file in the form of:

```toml
lower_snake_case_var_name = 'value' # in single quotes
```

### Configure the Services and Sequencer Networks with `networks-config.toml`

The `networks-config.toml` can be found in the `~/.astria/<instance>` directory.
It contains all of the configuration settings for which services are run by the
cli as well as a short list of useful overrides for different sequencer
networks.
Each network in the file will consist of a header section of the form
`[network.<network name>]`:

```toml
[networks.local]
```

And then a list of services that the cli will run and display in the TUI for
that network in the form of `[networks.<network name>.services.<service name>]`

```toml
[networks.local.services]
[networks.local.services.conductor]
# ...

[networks.local.services.conductor]
# ...
```

#### Choosing Specific Services

If you need to use a specific versions of a service that is not one of the
presets in the `networks-config.toml`, you can update that service for a given
network by updating the `version`, `download_url`, and `local_path` in the
services listed under each network.
For example, if you wanted to roll Composer back a version, you can do the
following:

Find the Composer within the `local` (or any other) network

```toml
# presets shown here
[networks.local.services.composer]
name = 'astria-composer'
version = 'v0.8.1'
download_url = 'https://github.com/astriaorg/astria/releases/download/composer-v0.8.1/astria-composer-aarch64-apple-darwin.tar.gz'
local_path = '<your home directory>/.astria/default/bin/astria-composer-v0.8.1'
```

And update it to a previous (or any other) version:

```toml
# update to a different version
[networks.local.services.composer]
name = 'astria-composer'
version = 'v0.7.0'
download_url = 'https://github.com/astriaorg/astria/releases/download/composer-v0.7.0/astria-composer-aarch64-apple-darwin.tar.gz'
local_path = '<your home directory>/.astria/default/bin/astria-composer-v0.7.0'
```

Then re-run `astria-go dev init` to download that service.

> NOTE: The 'name' and 'version' in the service config are used to name the
> binary that is downloaded from the 'download_url'. All downloaded binaries
> will be placed in the `~/.astria/<instance>/bin` directory and will be of the
> form `name-version` as specified in the service config within the
> `networks-config.toml` file.

The releases for all Astria services can be found [here](https://github.com/astriaorg/astria/releases).

If you are developing a service, you can also update the config to point at the
binary you are building locally:

```toml
# point the cli to your locally built binary
[networks.local.services.composer]
name = 'astria-composer'
version = 'dev'
download_url = ''
local_path = '<path to your local binary>'
```

#### Adding a New Network

If you would like to add a new network you can add a new section at the end of
the `networks-config.toml` file. For example, if you wanted to only run a local
sequencer (Astria Sequencer and CometBFT) you could add the following to the end
of the `networks-config.toml`:

```toml
[networks.sequencer_only]
sequencer_chain_id = 'sequencer-only'
sequencer_grpc = 'http://127.0.0.1:8080'
sequencer_rpc = 'http://127.0.0.1:26657'
rollup_name = 'astria-test-chain'
default_denom = 'nria'

[networks.sequencer_only.services]
[networks.sequencer_only.services.cometbft]
name = 'cometbft'
version = 'v0.38.8'
download_url = 'https://github.com/cometbft/cometbft/releases/download/v0.38.8/cometbft_0.38.8_darwin_arm64.tar.gz'
# make sure to update the path to the binary
local_path = '<your home directory>/.astria/default/bin/cometbft-v0.38.8'

[networks.sequencer_only.services.sequencer]
name = 'astria-sequencer'
version = 'v0.15.0'
download_url = 'https://github.com/astriaorg/astria/releases/download/sequencer-v0.15.0/astria-sequencer-aarch64-apple-darwin.tar.gz'
# make sure to update the path to the binary
local_path = '<your home directory>/.astria/default/bin/astria-sequencer-v0.15.0'
```

Then run your new network with the following commands:

```bash
astria-go dev init
astria-go dev run --network sequencer_only
```

> NOTE: It is best practice to re-run `dev init` after updating your
> `networks-config.toml`. This makes sure that all of the binaries that are used
> in all your network configurations have been downloaded prior to running.

### Interact with a Sequencer Network using the `sequencer-networks-config.toml`

The `sequencer-networks-config.toml` can be found in the top level `~/.astria`
directory. It provides presets for interacting with different sequencer networks
when using `astria-go sequencer` commands. You can use the `--network <network
name>` flag to simplify your commands when interacting with a sequencer network.
The config overrides the following flags when using the sequencer commands:

* `--sequencer-url`
* `--sequencer-chain-id`
* `--asset`
* `--fee-asset`

If you need to add presets for a new network, you can add a section to the end
of the `sequencer-networks-config.toml` file:

```toml
[networks.new_network]
sequencer_chain_id = 'new-network'
sequencer_url = '<rpc endpoint for the sequencer>'
asset = '<new asset>'
fee_asset = '<new fee asset>'
```

Then use the new config with:

```bash
astria-go sequencer nonce <other args> --network new_network
```

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

### TUI Logs

Because the TUI that is launched when using `dev run` manipulates the terminal,
the logs that would usually go to `stdout` or `stderr` are written to a log
file. You can find this log file in the `~/.astria/<instance>` directory that
the CLI is using for `dev run`.
For example, if you run `astria-go dev run`, the log file for the
currently running TUI will be at `~/.astria/default/astria-go.log`.

> NOTE: This log file gets overwritten every time you run the TUI. If you need
> to keep these logs, you will need to manually rename the file after closing
> the TUI.

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
