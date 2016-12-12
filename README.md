# Golang Microservices Example using gRPC

### Pre-requisites

Docker https://docs.docker.com/engine/installation

Protobuf v3

    $ brew install protobuf

Protoc-gen libraries:

    $ go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

Clone the repository:

    $ git clone git@github.com:dilipgurung/golang-microservices.git

### Protobufs

To regenerate Protocol Buffer files, run:

    $ make pb

### Build the application

    $ make build

### Run the application

    $ make up

Visit the web page in the browser:

[http://localhost:8000](http://localhost:8000/)

### Stop the application

    $ make down
