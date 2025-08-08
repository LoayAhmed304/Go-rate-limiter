package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/LoayAhmed304/GO-rate-limiter/internal/api"
	"github.com/LoayAhmed304/GO-rate-limiter/internal/logic/algorithms"
	"github.com/LoayAhmed304/GO-rate-limiter/pkg/logger"
)

func main() {
	fileName := flag.String("f", "./configs/config.json", "Path to the configuration file")
	port := flag.String("p", ":4000", `Port to run the server on. e.g. ":4000"`)
	flag.Parse()

	logger.Init()

	err := algorithms.InitConfigs(*fileName)
	if err != nil {
		logger.LogError("Failed to set up the configurations: " + err.Error())
		os.Exit(1)
	}

	mux := api.SetUpRoutes()
	logger.LogInfo("Starting the server. Listening on port " + *port)

	http.ListenAndServe(*port, mux)
}
