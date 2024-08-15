# Astria Go CLI

The `astria-go` CLI simplifies local rollup development and minimizes
dependencies. It provides functionality to easily run the Astria stack and
interact with the Sequencer.

## Table of Contents

* [Installation](#installation)
* [Running the Astria Sequencer](#running-the-astria-sequencer)
* [Instances](#instances)
* [Configuration](#configuration)
* [Development](#development)

## Installation

See all releases [here](https://github.com/astriaorg/astria-cli-go/releases).

### Install From GitHub Release

1. Download the latest release for your platform:

   ```bash
   # For macOS silicon (arm64)
   export RELEASE_URL="https://github.com/astriaorg/astria-cli-go/releases/download/v0.12.0/astria-go-v0.12.0-darwin-arm64.tar.gz"
   curl -L $RELEASE_URL --output astria-go.tar.gz
   ```

2. Extract the binary:

   ```bash
   tar -xzvf astria-go.tar.gz
   ```

3. Verify the installation:

   ```bash
   ./astria-go version
   ```

4. Optionally, move the binary to a location in your PATH:

   ```bash
   mv astria-go /usr/local/bin/
   ```

### Install Nightly Release

1. Download the nightly release:

   ```bash
   export NIGHTLY_URL="download url of the build you need"
   curl -L $NIGHTLY_URL > astria-cli.tar.gz
   ```

2. Extract and install:

   ```bash
   tar -xvzf astria-cli.tar.gz
   mv astria-go /usr/local/bin/
   ```

3. Verify the installation:

   ```bash
   astria-go version
   ```

### Build Locally from Source

Prerequisites:

* [GO](https://go.dev/doc/install)
* [just](https://github.com/casey/just)

Steps:

1. Clone the repository:

   ```bash
   git clone git@github.com:astriaorg/astria-cli-go.git
   cd astria-cli-go
   ```

2. Build the project:

   ```bash
   just build
   ```

3. Verify the build:

   ```bash
   just run "version"
   # or
   go run main.go version
   ```

## Running the Astria Sequencer

### Initialize Configuration

```bash
astria-go dev init
```

This command downloads binaries, generates configuration files, and initializes
CometBFT. Files are organized in the `~/.astria/<instance>` directory.

### Usage

#### Run a Local Sequencer

```bash
astria-go dev run --network local
```

This command starts a local Sequencer (Cometbft and Astria-Sequencer),
Conductor, and Composer using pre-built binaries.

#### Run Against a Remote Sequencer

To run Composer and Conductor locally against a remote Astria Sequencer:

```bash
# Run against the Astria Dusk dev net 
astria-go dev run --network dusk
```

When using a remote sequencer, you'll need to create an account on the remote
sequencer. For more details, refer to the [Astria
documentation](https://docs.astria.org/developer/tutorials/run-local-rollup-against-remote-sequencer#configure-the-local-astria-components).

### Run Custom Binaries

You can run components from a local monorepo during development. For example, to
use a locally compiled Conductor:

```bash
astria-go dev run --network local \
  --conductor-path <absolute path to the Astria mono repo>/target/debug/astria-conductor
```

Or update the local path in the `~/.astria/<instance>/networks-config.toml`:

```toml
[networks.local.services.conductor]
name = 'astria-conductor'
version = 'dev'
download_url = ''
local_path = 'path to your local conductor'
args = []
```

You can swap out some or all binaries:

```bash
astria-go dev run --network local \
  --sequencer-path <sequencer bin path> \
  --cometbft-path <cometbft bin path> \
  --composer-path <composer bin path> \
  --conductor-path <conductor bin path>
```

### Interact with the Sequencer

Use `astria-go sequencer [command]` to interact with the sequencer. For a full
list of commands:

```bash
astria-go sequencer -h
```

Use the `--network` flag to configure which sequencer network the commands will
run against.

## Instances

Use the `--instance` flag to manage multiple rollups:

```bash
astria-go dev init
astria-go dev init --instance hello
astria-go dev init --instance world
```

This creates separate directories in `~/.astria/` for each instance, containing
configs and binaries for running the Astria stack.

## Configuration

The CLI uses three configuration files:

1. `base-config.toml`: Sets service environment variables
2. `networks-config.toml`: Configures services and sequencer networks
3. `sequencer-networks-config.toml`: Used for `astria-go sequencer` commands

### Set Service Environment Variables

Edit `~/.astria/<instance>/config/base-config.toml` to add or change settings:

```toml
lower_snake_case_var_name = 'value'
```

### Configure Networks and Services

The `~/.astria/<instance>/networks-config.toml` file configures which services
are run and provides overrides for different sequencer networks.

#### Customizing Service Versions

To use a specific version of a service, update the `version`, `download_url`,
and `local_path` in the `networks-config.toml`. For example, to roll back
Composer:

```toml
[networks.local.services.composer]
name = 'astria-composer'
version = 'v0.7.0'
download_url = 'https://github.com/astriaorg/astria/releases/download/composer-v0.7.0/astria-composer-aarch64-apple-darwin.tar.gz'
local_path = '<your home directory>/.astria/default/bin/astria-composer-v0.7.0'
args = []
```

Then run `astria-go dev init` to download the specified service.

For local development, point to your locally built binary:

```toml
[networks.local.services.composer]
name = 'astria-composer'
version = 'dev'
download_url = ''
local_path = '<path to your local binary>'
args = []
```

#### Run a Generic Service

Add a local service to your network:

```toml
[networks.local.services.echo]
name = 'echo'
version = ''
download_url = ''
local_path = '/bin/bash'
args = ['-c', 'echo -e "hello world\nhello again!"']
```

> NOTE: All arguments are interpreted literally. This may affect how some
> service arguments are parsed. Running through a bash shell command
> resolves this issue in the case of running an `echo` command above.

Add a service from a release:

```toml
[networks.local.services.your_service]
name = 'your_service'
version = 'v0.0.0'
download_url = 'download url to the release'
local_path = '<your home directory>/.astria/default/bin/<your_service_name-version>'
args = ['your', 'service', 'args']
```

#### Adding a New Network

To add a new network, append a new section to `networks-config.toml`:

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
local_path = '<your home directory>/.astria/default/bin/cometbft-v0.38.8'
args = []

[networks.sequencer_only.services.sequencer]
name = 'astria-sequencer'
version = 'v0.15.0'
download_url = 'https://github.com/astriaorg/astria/releases/download/sequencer-v0.15.0/astria-sequencer-aarch64-apple-darwin.tar.gz'
local_path = '<your home directory>/.astria/default/bin/astria-sequencer-v0.15.0'
args = []
```

Run your new network with:

```bash
astria-go dev init
astria-go dev run --network sequencer_only
```

### Interact with Sequencer Networks

The `~/.astria/sequencer-networks-config.toml` provides presets for interacting
with different sequencer networks when using `astria-go sequencer` commands. Use
the `--network <network name>` flag to simplify your commands.

To add presets for a new network, append to `sequencer-networks-config.toml`:

```toml
[networks.new_network]
sequencer_chain_id = 'new-network'
sequencer_url = '<rpc endpoint for the sequencer>'
asset = '<new asset>'
fee_asset = '<new fee asset>'
```

Use the new config with:

```bash
astria-go sequencer nonce <other args> --network new_network
```

## Development

Requirements:

* Go version 1.22 or newer
* Updated `gopls` settings for correct parsing of build tags

The CLI uses [Cobra](https://github.com/spf13/cobra) for structuring the
command-line interface. To add new commands:

1. Install Cobra CLI:

   ```bash
   go install github.com/spf13/cobra-cli@latest
   ```

2. Add a new command:

   ```bash
   cobra-cli add <command-name>
   ```

### TUI Logs

TUI logs are written to `~/.astria/<instance>/astria-go.log`. This file is
overwritten each time you run the TUI.

### Testing

#### Unit Tests

```bash
# Run unit tests
just test
```

#### Integration Tests

The integration tests require a running local sequencer network to test against.
Prior to testing, run:

```bash
just run dev run --network local
```

In a new terminal window run:

```bash
# Run integration tests
just test-integration

# Run all tests
just test-all
```
