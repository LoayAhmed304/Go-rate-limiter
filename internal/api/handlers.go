package api

import (
	"net/http"

	"github.com/LoayAhmed304/GO-rate-limiter/internal/logic/algorithms"
	"github.com/LoayAhmed304/GO-rate-limiter/pkg/logger"
)

func HandleRateLimit(w http.ResponseWriter, r *http.Request) {

	key, route := r.Header.Get("X-Rate-Limit-Key"), r.Header.Get("X-Original-Path")

	valid, timeLeft := algorithms.ConfigInstance.Algorithm.AllowRequest(key, route)

	if valid {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte("Request allowed"))
		if err != nil {
			logger.LogError("Failed to write response: " + err.Error())
		}
	} else {
		w.Header().Add("Retry-After", timeLeft.String())
		w.WriteHeader(http.StatusTooManyRequests)

		_, err := w.Write([]byte("Too many requests. Try again later."))
		if err != nil {
			logger.LogError("Failed to write response: " + err.Error())
		}
	}

}

func ValidateHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if key := r.Header.Get("X-Rate-Limit-Key"); key == "" {
			w.WriteHeader(http.StatusBadRequest)

			_, err := w.Write([]byte("Client key not found"))
			if err != nil {
				logger.LogError("Failed to write response: " + err.Error())
				return
			}
		}

		if route := r.Header.Get("X-Original-Path"); route == "" {
			r.Header.Set("X-Original-Path", r.URL.Path)
		}

		next.ServeHTTP(w, r)
	})

}
