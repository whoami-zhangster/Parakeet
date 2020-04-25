package main

import (
	"time"

	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	conf "github.com/whoami-zhangster/Parakeet/pkg/config"
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

	// create endpoint w/ operations and payload
	r := mux.NewRouter()
	r.HandleFunc(endpointConfig.Path, endpointConfig.CreateHandleFunc())

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
