package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/whoami-zhangster/Parakeet/pkg/logger"
)

type (
	HttpServerConfig struct {
		API            map[string]APIConfig `yaml:"api"` // path to api config
		Port           int                  `yaml:"port"`
		ReadTimeout    int                  `yaml:"readTimeout"`
		WriteTimeout   int                  `yaml:"writeTimeout"`
		MaxHeaderBytes int                  `yaml:"maxHeaderBytes"`
	}

	HttpServer struct {
		log            logger.Logger
		API            map[string]*API `yaml:"api"`
		serverMap      map[string]*http.Server
		Port           int `yaml:"port"`
		ReadTimeout    int `yaml:"readTimeout"`
		WriteTimeout   int `yaml:"writeTimeout"`
		MaxHeaderBytes int `yaml:"maxHeaderBytes"`
	}

	// API is a rest API definition
	APIConfig struct {
		Methods map[string]ResponseConfig `yaml:"methods"`
	}

	API struct {
		Methods map[string]ResponseConfig // method to response
		log     logger.Logger
		path    string
		port    int
		srv     *http.Server
	}

	ResponseConfig struct {
		Status     string            `yaml:"status"`     // e.g. "200 OK"
		StatusCode int               `yaml:"statusCode"` // e.g. 200
		Header     map[string]string `yaml:"headers"`    // kv rep
		Body       string            `yaml:"body"`
	}
)

func NewHttpServer(log logger.Logger, config HttpServerConfig) *HttpServer {
	m := make(map[string]*API)
	for path, conf := range config.API {
		m[path] = newAPI(log, path, config.Port, conf)
	}
	return &HttpServer{
		log:            log,
		Port:           config.Port,
		API:            m,
		serverMap:      make(map[string]*http.Server),
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
}

func newAPI(log logger.Logger, path string, port int, config APIConfig) *API {
	return &API{
		Methods: config.Methods,
		log:     log,
		path:    path,
		port:    port,
	}
}

func (hs *HttpServer) CreateServers() {
	for _, api := range hs.API {
		// Run ea. server in goroutine
		api.CreateAndRunServer()
	}
}

func (api *API) CreateAndRunServer() {
	api.log.Infof("Creating server for api: %+v", api)
	router := api.createRouter(api.path)
	address := fmt.Sprintf("0.0.0.0:%d", api.port)
	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// assign created server to api
	api.srv = srv

	err := srv.ListenAndServe()
	if err != nil {
		api.log.Error(err)
	}

	// Run server in goroutine
	/*
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				api.log.Error((err))
			}
		}()
	*/
}

func (api *API) KillServer() error {
	if api.srv == nil {
		return errors.New("no server to kill")
	}

	return api.srv.Close()
}

func (api *API) createRouter(path string) *mux.Router {
	r := mux.NewRouter()
	for method, responseConfig := range api.Methods {
		api.log.Infof("creating router for path %s method %s response conifg %v", path, method, responseConfig)
		r.HandleFunc(path, responseConfig.CreateHandleFunc()).
			Methods(method).
			Schemes("http")
	}
	return r
}

// CreateHandleFunc Creates a http.Handler function from an endpoint definition
func (resp ResponseConfig) CreateHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp.copyValues(w.Header())
		w.WriteHeader(resp.StatusCode)
		fmt.Fprintf(w, resp.Body)
	}
}

func (resp ResponseConfig) copyValues(header http.Header) {
	for k, v := range resp.Header {
		header.Set(k, v)
	}
}

func stringIsJSON(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil

}
