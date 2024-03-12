default:
  @just --list

build:
    go build -o bin/astria-dev 

fmt:
    go fmt ./...

clean:
    rm -rf bin
    rm -rf local-dev-astria
    rm -rf data

init:
    ./bin/astria-dev dev init

run: 
    ./bin/astria-dev dev run
