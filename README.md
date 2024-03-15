# The Astria Development CLI

The `astria-dev` cli is a tool designed to make local rollup development as
simple and dependency free as possible.

Currently the cli only works on arm Macos.

# Locally Build the CLI

Dependencies: (only required for development)

- [GO](https://go.dev/doc/install)
- [just](https://github.com/casey/just)

```
git clone git@github.com:astriaorg/astria-cli-go.git
cd astria-cli-go
git checkout feat/local-binaries
just build
just run "dev init"
just run "dev run"
```

This will download, configure, and run the following binaries of these applications:

| App              | Version |
| ---------------- | ------- |
| Cometbft         | v0.37.4 |
| Astria-Sequencer | v0.9.0  |
| Astria-Conductor | v0.12.0 |
| Astria-Composer  | v0.4.0  |

And place them in a `local-dev-astria` directory, along with several other
configuration files for everything.

The cli runs the minimum viable components for testing a rollup against the
Astria stack, allowing developers to confirm that their rollup interacts with
Astria's apis correctly.

## Development

Requires go version 1.20 or newer.

You may also need to update your `gopls` settings in your editor for build tags to allow for
correct parsing of the build tags in the code. This will depend on your IDE, but
for VS Code you can open your settings and add:

```
"gopls": {
    "buildFlags": ["-tags=darwin arm64 amd64 linux"]
}
```

The cli is built using [Cobra](https://github.com/spf13/cobra). Once you've
pulled the repo you can install the `cobra-cli` as follows to add new commands
for development:

```
# install cobra-cli
go install github.com/spf13/cobra-cli@latest
# add new command, e.g. `run`
cobra-cli add <new-command>
```

### Available Commands

| Command         | Description                                                          |
| --------------- | -------------------------------------------------------------------- |
| `version`       | Prints the cli version.                                              |
| `help`          | Show help.                                                           |
| `dev`           | Root command development.                                            |
| `dev init`      | Downloads binaries and initializes the local environment.            |
| `dev run`       | Runs a minimal, local Astria stack.                                  |
| `dev clean`     | Deletes the local data for the Astria stack.                         |
| `dev clean all` | Deletes the local data and downloaded binaries for the Astria stack. |
