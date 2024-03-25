default:
  @just --list

build:
    go build -o bin/astria-dev 

# format all go files
fmt:
    go fmt ./...

defaultargs := ''
run args=defaultargs:
    go run main.go {{args}} > tview_log.txt 2>&1

# kill all processes that may be hanging because of improper shutdown of the tview app
pskill:
    ps aux | grep -E '[c]omposer|[c]onductor|[s]equencer|[c]ometbft' | awk '{print $2}' | xargs kill -9
