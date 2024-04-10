# list all available commands
default:
    @just --list

# build the binary for the cli
build:
    go build -o bin/astria-go
alias b := build

# test go code
test:
    go test ./...
alias t := test

# format all go files
fmt:
    go fmt ./...
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
    golangci-lint run

[no-exit-message]
_lint-md:
    markdownlint-cli2 "**/*.md" "#bin" "#.github"

defaultargs := ''
# run the cli. takes quoted cli command to run, e.g. `just run "dev init"`. logs cli output to tview_log.txt
run args=defaultargs:
    go run main.go {{args}}
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
