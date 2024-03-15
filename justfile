default:
  @just --list

build:
    go build -o bin/astria-dev 

fmt:
    go fmt ./...

defaultargs := ''
run args=defaultargs:
    go run main.go {{args}}
