build:
    go build -o bin/astria-local 

clean:
    rm -rf bin
    rm -rf local-dev-astria

run:
    ./bin/astria-local init
