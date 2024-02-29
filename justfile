build:
    go build -o bin/astria-dev 

clean:
    rm -rf bin
    rm -rf local-dev-astria
    rm -rf data

init:
    ./bin/astria-dev init

run: 
    ./bin/astria-dev run
