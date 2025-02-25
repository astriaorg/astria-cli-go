# list all available commands
default:
    @just --list

default_binary_name := 'astria-go'

# build the binary for the cli
build-cli binary_name=default_binary_name:
    cd modules/cli && go build -o ../../bin/{{binary_name}}
alias b := build-cli

install-cli binary_name=default_binary_name:
    @just build-cli {{binary_name}}
    mv bin/{{binary_name}} ~/go/bin/{{binary_name}}

# test go code, minus integration tests
test:
    cd modules/cli && go test ./...
    cd modules/go-sequencer-client && go test ./...
    cd modules/bech32m && go test ./...
alias t := test

# unit tests with coverage report that opens in browser
test-cov:
    cd modules/cli && go test ./... -coverprofile=coverage.out
    cd modules/cli && go tool cover -html=coverage.out
    cd modules/go-sequencer-client && go test ./... -coverprofile=coverage.out
    cd modules/go-sequencer-client && go tool cover -html=coverage.out
    cd modules/bech32m && go test ./... -coverprofile=coverage.out
    cd modules/mech32m && go tool cover -html=coverage.out

# run unit and integration tests, and tests that require tty.
test-all: test test-integration-cli
alias ta := test-all
# run integration tests for all modules
test-integration: test-integration-cli test-integration-go-sequencer-client
alias ti := test-integration

# run integrations tests. requires running geth + cometbft + astria core.
test-integration-cli:
    just build-cli astria-go-testy
    cd modules/cli/integration_tests && go test -v ./... -tags=integration_tests -count=1
    rm ./bin/astria-go-testy

# run integration tests for go-sequencer-client. requires running geth + cometbft + astria core.
test-integration-go-sequencer-client:
    cd modules/go-sequencer-client/integration_tests && go test ./... -tags=integration_tests -count=1

cleanup-integration-tests:
    rm -f ./bin/astria-go-testy

# format all go files
fmt:
    cd modules/cli && go fmt ./...
    cd modules/go-sequencer-client && go fmt ./...
    cd modules/bech32m && go fmt ./...
alias f := fmt

default_lang := 'all'

# Can lint 'go', 'md', or 'all'. Defaults to all.
lint lang=default_lang:
    @just _lint-{{lang}}
alias l := lint

@_lint-all:
    @just _lint-go
    @just _lint-md

[no-exit-message]
_lint-go:
    cd modules/cli && golangci-lint run
    cd modules/go-sequencer-client && golangci-lint run
    cd modules/bech32m && golangci-lint run

[no-exit-message]
_lint-md:
    markdownlint-cli2 "**/*.md" "#bin" "#.github" --config .markdownlint.json

# fix markdown linting issues that can be auto-fixed
fix-md:
    markdownlint-cli2 "**/*.md" "#bin" "#.github" --fix

defaultargs := ''
# build and run the cli with --log-level=debug. The process will be named astria-go-dev.
run *args=defaultargs:
    cd modules/cli && go build -o ../../bin/astria-go-dev && ../../bin/astria-go-dev {{args}} --log-level=debug

alias r := run

run-race args=defaultargs:
    go run -race main.go {{args}}

# show any running Astria processes
[no-exit-message]
@pscheck:
    ps aux | grep -E '[c]omposer|[c]onductor|[s]equencer|[c]ometbft'

# kill all processes that may be hanging because of improper shutdown of the tview app
pskill:
    ps aux | grep -E '[c]omposer|[c]onductor|[s]equencer|[c]ometbft' | awk '{print $2}' | xargs kill -9
