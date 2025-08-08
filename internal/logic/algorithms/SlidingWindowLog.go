package algorithms

import (
	"time"

	"github.com/LoayAhmed304/GO-rate-limiter/internal/configs"
	"github.com/LoayAhmed304/GO-rate-limiter/internal/logic/algorithms/structures"
)

var windowSize time.Duration
var maxRequests int

// AllowRequest determines whether a request may be allowed for a client,
// based on the Sliding Window Log algorithm.
//
// It takes a slice of timestamps representing the previous requests from the client,
// and returns whether the current request is allowed.
// If not allowed, it alos returns the time remaining until the next allowed request.
func AllowRequest(clientIP, route string) (bool, time.Duration) {
	// route is not assigned in configs, so we assume it is whitelisted
	if (*structures.ClientsLogs)[route] == nil {
		return true, 0
	}

	windowSize = configs.ConfigInstance.RoutesConfigs[route].Interval
	maxRequests = configs.ConfigInstance.RoutesConfigs[route].Limit

	_, exists := (*structures.ClientsLogs)[route][clientIP]
	if !exists {
		s := make([]time.Time, 0, maxRequests)
		(*structures.ClientsLogs)[route][clientIP] = &s
	}

	clientLogs := (*structures.ClientsLogs)[route][clientIP]

	if curRequests := len(*clientLogs); curRequests < maxRequests {
		*clientLogs = append(*clientLogs, time.Now())

		go clearLogs(clientLogs) // a goroutine to clear the expired timestamps

		return true, 0
	}

	clearLogs(clientLogs)

	if curRequests := len(*clientLogs); curRequests < maxRequests {
		*clientLogs = append(*clientLogs, time.Now())
		return true, 0
	}

	remainingTime := calcRemainingTime(clientLogs)

	return false, remainingTime
}

// clearLogs is a private function used as a helper function for the AllowRequest
// in SlidingWindowLog algorithm.
//
// It removes all outdated requests in the given slice of timestamps
func clearLogs(clientLogs *[]time.Time) {
	if len(*clientLogs) == 0 {
		return
	}

	curWindowStart := time.Now().Add(-windowSize)
	curTimestamp := (*clientLogs)[0]
	for len(*clientLogs) > 0 && curTimestamp.Before(curWindowStart) {
		*clientLogs = (*clientLogs)[1:]

		if len(*clientLogs) >= 1 {
			curTimestamp = (*clientLogs)[0]
		}
	}
}

func calcRemainingTime(clientLogs *[]time.Time) time.Duration {
	if len(*clientLogs) == 0 {
		return 0
	}

	windowStartTime := time.Now().Add(-windowSize)
	firstRequestTime := (*clientLogs)[0]
	remainingTime := firstRequestTime.Sub(windowStartTime)

	return remainingTime
}
