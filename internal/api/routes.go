package api

import "net/http"

func SetUpRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", HandleRateLimit)

	return ValidateHeaders(mux)
}
