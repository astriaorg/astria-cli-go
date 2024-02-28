build:
    go build -o bin/astria-local 

clean:
    rm -rf bin
    rm -rf local-dev-astria
    rm -rf data

init:
    ./bin/astria-local init
