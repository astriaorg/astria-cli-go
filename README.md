# Locally Build the CLI

Dependencies: (for now)

- GO
- just
- mprocs

```
git clone git@github.com:astriaorg/astria-cli-go.git
cd astria-cli-go
git checkout feat/local-binaries
just build
just init
just run
```

`just init` is an alias for `./bin/astria-dev init`
`just run` is an alias for `./bin/astria-dev run`

This will download, configure, and run the following binaries of these applications:

| App       | Version |
| --------- | ------- |
| Cometbft  | v0.37.4 |
| Sequencer | v0.9.0  |
| Conductor | v0.12.0 |
| Composer  | v0.4.0  |

And place them in a `local-dev-astria` directory, along with several other
configuration files for everything.

### Useful Commands While Testing

```
# In astria-cli-go/
# remove locally built cli binary, the local-dev-astria/ dir, and the data/ dir
just clean
```
