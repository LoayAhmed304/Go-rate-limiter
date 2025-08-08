package main

import (
	"flag"
	"os"

	"github.com/LoayAhmed304/GO-rate-limiter/internal/configs"
	"github.com/LoayAhmed304/GO-rate-limiter/pkg/logger"
)

func main() {
	fileName := flag.String("f", "./configs/config.json", "Path to the configuration file")
	flag.Parse()

	logger.Init()

	err := configs.ParseConfig(*fileName)
	if err != nil {
		logger.LogError("Failed to set up the configurations: " + err.Error())
		os.Exit(1)
	}
}
