package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	conf "github.com/whoami-zhangster/Parakeet/pkg/config"
	"github.com/whoami-zhangster/Parakeet/pkg/rest"
)

func main() {
	// create logger
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()
	log.Info("created logger")

	// Config file from command line args as byte array
	endpointConfig, err := conf.Setup(log)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	srv := rest.NewHttpServer(log, *endpointConfig)

	srv.CreateServers()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	// kill all
	srv.KillAllAPI()

}
