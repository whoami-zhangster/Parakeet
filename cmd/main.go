package main

import (
	"fmt"

	conf "github.com/whoami-zhangster/Parakeet/pkg/config"
)

func main() {
	// Config file from command line args as byte array
	endpointConfig, err := conf.Setup()
	if err != nil {
		panic("TODO ")
	}

	fmt.Printf("\n%+v\n", *endpointConfig)

	// create endpoint w/ operations and payload
}
