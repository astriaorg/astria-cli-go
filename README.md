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
just init
just run
```

`just init` is an alias for `./bin/astria-dev dev init`
`just run` is an alias for `./bin/astria-dev dev run`

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

Requires go version 1.17 or newer.

You may also need to update your `gopls` settings in your editor for build tags to allow for
correct parsing of the build tags in the code. This will depend on your IDE, but
for VS Code you can open your settings and add:

```
"gopls": {
    "buildFlags": ["-tags=darwin arm64 amd64 linux"]
}
```

### Useful Commands While Testing

```
# In astria-cli-go/
# removes the locally built `astria-dev` cli binary, the local-dev-astria/ dir, and the data/ dir
just clean
```
