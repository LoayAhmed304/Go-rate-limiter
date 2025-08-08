package api

import (
	"net/http"

	"github.com/LoayAhmed304/GO-rate-limiter/internal/logic/algorithms"
	"github.com/LoayAhmed304/GO-rate-limiter/pkg/logger"
)

func HandleRateLimit(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr

	if clientIP == "" {
		w.WriteHeader(http.StatusBadRequest)

		_, err := w.Write([]byte("Client IP not found"))
		if err != nil {
			logger.LogError("Failed to write response: " + err.Error())
		}

		return
	}

	route := r.URL.Path

	valid, timeLeft := algorithms.ConfigInstance.Algorithm.AllowRequest(clientIP, route)

	if valid {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte("Request allowed"))
		if err != nil {
			logger.LogError("Failed to write response: " + err.Error())
		}
	} else {
		w.WriteHeader(http.StatusTooManyRequests)

		_, err := w.Write([]byte("Too many requests. Try again in " + timeLeft.String()))
		if err != nil {
			logger.LogError("Failed to write response: " + err.Error())
		}
	}

}
