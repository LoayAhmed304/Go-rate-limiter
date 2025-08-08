package configs

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"

	"github.com/LoayAhmed304/GO-rate-limiter/pkg/logger"
)

type RouteConfig struct {
	Limit    int           `json:"limit"`
	Interval time.Duration `json:"interval"`
}

type Config struct {
	Algorithm     string                 `json:"algorithm"`
	RoutesConfigs map[string]RouteConfig `json:"routes"`
}

type rawConfig struct {
	Algorithm string           `json:"algorithm"`
	Routes    []rawRouteConfig `json:"routes"`
}

type rawRouteConfig struct {
	Route    string `json:"route"`
	Limit    int    `json:"limit"`
	Interval string `json:"interval"`
}

var ConfigInstance *Config = &Config{}

// ParseConfig parses the given json configuration file
// and populate the ConfigInstance struct with the necessary data.
func ParseConfig(fileName string) error {
	if extension := path.Ext(fileName); extension != ".json" {
		return errors.New("invalid configuration file format: " + extension + ". Expected .json")
	}

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	var rawConfigInstance rawConfig // temp variable to hold the raw configuration data
	err = json.NewDecoder(file).Decode(&rawConfigInstance)
	if err != nil {
		return err
	}

	err = convertRawConfig(&rawConfigInstance)
	if err != nil {
		return err
	}

	logger.LogInfo("Configurations set up successfully")
	return nil
}

// convertRawConfig converts the raw configuration structure into a mapped structure
// that is used in the application for faster access and better performance.
//
// It also converts the interval string into a time.Duration type, that's the only error it might return
// in case of failure.
func convertRawConfig(raw *rawConfig) error {
	ConfigInstance.Algorithm = raw.Algorithm
	ConfigInstance.RoutesConfigs = make(map[string]RouteConfig, len(raw.Routes))
	for _, route := range raw.Routes {
		interval, err := time.ParseDuration(route.Interval)
		if err != nil {
			return err
		}

		ConfigInstance.RoutesConfigs[route.Route] = RouteConfig{
			Limit:    route.Limit,
			Interval: interval,
		}
	}
	return nil
}
