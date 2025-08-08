package api

import (
	"net/http"

	"github.com/LoayAhmed304/GO-rate-limiter/internal/logic/algorithms"
)

func HandleRateLimit(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr

	if clientIP == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Client IP not found"))
		return
	}

	route := r.URL.Path

	valid, timeLeft := algorithms.ConfigInstance.Algorithm.AllowRequest(clientIP, route)

	if valid {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request allowed"))
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many requests. Try again in " + timeLeft.String()))
	}

}
