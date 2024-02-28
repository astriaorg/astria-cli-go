# Locally Build the CLI

Dependencies: (for now)

- GO
- just
- mprocs

```
git clone git@github.com:astriaorg/astria-cli-go.git
cd astria-cli-go
just build
just init
```

`just init` is an alias for `./bin/astria-local init`

This will download the following binaries of these applications:

| App       | Version |
| --------- | ------- |
| Cometbft  | v0.37.4 |
| Sequencer | v0.9.0  |
| Conductor | v0.12.0 |
| Composer  | v0.4.0  |

And place them in a `local-dev-astria` directory, along with several other
configuration files for everything.

# Run the Applications

NOTE: this will eventually be integrated into the cli as `astria-local run --local`

```
cd local-dev-astria
mprocs
```

### Useful Commands While Testing

```
# In astria-cli-go/
# remove locally built cli binary and the local-dev-astria/ dir
just clean

# In astria-cli-go/local-dev-astria/
# remove local data for cometbft and sequencer
just clean
```
