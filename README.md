# storage

storage provides gRPC server and client for image uploading and storage.

## Installation

`go get github.com/buraksekili/storage` which adds this repository into 
`$GOPATH/src/github.com/buraksekili/storage`

or, you can directly clone this repository via;

`git clone https://github.com/buraksekili/storage.git`

### Usage

`make server` => starts the server

`make client` => sends an image to the server

### Flags

```
-path   (string)    defines the path of the image

-addr   (string)    defines the address of the server, e.g, 'localhost'

-port   (string)    defines the port for the server
```

`-addr` and `-port` Flags are available for both server and client.


__`Caveat`__: The port and address of the client and server must be same.

### Usage

```bash
~ go run cmd/client/main.go -path="./img/test-client/f.png"
```

## Requirements

`protoc` is required to generate new pb files. 
Therefore, it is a must for running commands indicated in the Makefile.

