package config

import (
	"errors"
	"flag"
	"io/ioutil"

	"github.com/whoami-zhangster/Parakeet/pkg/logger"
	"github.com/whoami-zhangster/Parakeet/pkg/rest"

	"github.com/go-yaml/yaml"
)

// Setup takes the configFile path as a command line argument, loads the file and parses it into Endpoint struct
func Setup(log logger.Logger) (*rest.HttpServerConfig, error) {
	var configFile string
	flag.StringVar(&configFile, "config", "config.yaml", "config filepath for yaml file")
	flag.Parse()
	if configFile == "" {
		return nil, errors.New("no config file provided")
	}

	log.Debugf("reading file %s", configFile)

	cf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	endpoint := &rest.HttpServerConfig{}

	err = parse(cf, endpoint)
	if err != nil {
		return nil, err
	}

	return endpoint, nil

}

func parse(yamlFileInBytes []byte, dist interface{}) error {
	if dist == nil {
		return errors.New("can not cast to nil config")
	}
	err := yaml.Unmarshal(yamlFileInBytes, dist)
	if err != nil {
		return err
	}
	return nil
}
