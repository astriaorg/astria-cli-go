build:
    go build -o bin/astria-local 

clean:
    rm -rf bin
    rm -rf local-dev-astria

init:
    ./bin/astria-local init
