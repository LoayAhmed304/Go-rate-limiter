package algorithms

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"

	"github.com/LoayAhmed304/GO-rate-limiter/pkg/logger"
)

var ConfigInstance *Config

// InitConfigs parses the given json configuration file
// and populate the ConfigInstance struct with the necessary data.
//
// It first parses the json file into a rawConfig structure,
// then converts it into a global ConfigInstance structure for easier acces.
//
// Then it initializes a global algorithm interface based on the configuration,
// and finally initializes the ClientsLogs map for every route defined in the configurations.
func InitConfigs(fileName string) error {
	rawConfigInstance := &rawConfig{}

	err := parseConfigFile(fileName, rawConfigInstance)
	if err != nil {
		return err
	}

	err = fillConfigInstance(rawConfigInstance)
	if err != nil {
		return err
	}

	InitAlgorithm(rawConfigInstance.Algorithm)

	routes := getMapKeys(ConfigInstance.RoutesConfigs)
	ConfigInstance.Algorithm.Init(routes)

	logger.LogInfo("Configurations set up successfully")
	return nil
}

// fillConfigInstance converts the raw configuration structure into a mapped structure
// that is used in the application for faster access and better performance.
//
// It also converts the interval string into a time.Duration type, that's the only error it might return
// in case of failure.
func fillConfigInstance(raw *rawConfig) error {
	ConfigInstance = &Config{}
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

func getMapKeys(m map[string]RouteConfig) []string {
	keys := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}

// InitAlgorithm initializes the global Algorithm interface based on the given algorithm name.
//
// If the algorithm is not recognized, it defaults to SlidingWindowLog.
func InitAlgorithm(algorithm string) {
	switch algorithm {
	case "TokenBucket":
		ConfigInstance.Algorithm = &TokenBucket{}
	case "SlidingWindowLog":
		ConfigInstance.Algorithm = &SlidingWindowLog{}
	default:
		ConfigInstance.Algorithm = &SlidingWindowLog{}
	}
}

// parseConfigFile reads the configuration file from the given path and
// decodes it into a given rawConfig structure
func parseConfigFile(fileName string, rawConfigInstance *rawConfig) error {
	if extension := path.Ext(fileName); extension != ".json" {
		return errors.New("invalid configuration file format: " + extension + ". Expected .json")
	}

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(rawConfigInstance)
	if err != nil {
		return err
	}
	return nil
}
