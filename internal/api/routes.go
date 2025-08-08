package api

import "net/http"

func SetUpRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", HandleRateLimit)

	return mux
}
