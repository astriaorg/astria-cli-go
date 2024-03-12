default:
  @just --list

build:
    go build -o bin/astria-dev 

fmt:
    go fmt ./...

clean:
    ./bin/astria-dev dev clean

cleanall:
    ./bin/astria-dev dev clean all

init:
    ./bin/astria-dev dev init

run: 
    ./bin/astria-dev dev run
