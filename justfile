# list all available commands
default:
  @just --list

# build the binary for the cli
build:
    go build -o bin/astria-cli

# test go code
test:
    go test ./...

alias t := test

# format all go files
fmt:
    go fmt ./...

defaultargs := ''
# run the cli. takes quoted cli command to run, e.g. `just run "dev init"`. logs cli output to tview_log.txt
run args=defaultargs:
    go run main.go {{args}} > tview_log.txt 2>&1

# kill all processes that may be hanging because of improper shutdown of the tview app
pskill:
    ps aux | grep -E '[c]omposer|[c]onductor|[s]equencer|[c]ometbft' | awk '{print $2}' | xargs kill -9
