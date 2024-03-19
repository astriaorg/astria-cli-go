default:
  @just --list

build:
    go build -o bin/astria-dev 

fmt:
    go fmt ./...

defaultargs := ''
run args=defaultargs:
    go run main.go {{args}} > tview_log.txt 2>&1

pskill:
    ps aux | grep -E '[c]omposer|[c]onductor|[s]equencer|[c]ometbft|[g]eth' | awk '{print $2}' | xargs kill -9
