# Parakeet: a go library for creating virtual dummy endpoints

## Project Features
The goal of this project is to provide an easy to use, extensible Mock API server that can be used for local mock integration testing, performance testing, etc

Current Features:
* Parsing of API configuration from yaml file ex in `example.yml`
* Creating REST HTTP server based on parsed specifications

Planned Features:
* Dynamic tearing down/bringing up of HTTP server
* Creating own REST service to allow users modify these a REST HTTP Server dynamically
* Rate Limiting/throtteling responses
* gRPC support
* define FSM to test basic actions interactively


## Getting Started
```console
go get github.com/whoami-zhangster/Parakeet
```
### Build and Run
```console
go build -o main cmd/main.go
./main -config <yaml file>
```