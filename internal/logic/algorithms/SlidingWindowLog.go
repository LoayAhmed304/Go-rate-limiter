package algorithms

import (
	"time"
)

// assume it's 5 requests per 10 seconds for now
// project structure will be adjusted to support given configs for each algorithm
var maxRequests int = 5
var windowSize int = 10 // 10 seconds

// AllowRequest determines whether a request may be allowed for a client,
// based on the Sliding Window Log algorithm.
//
// It takes a slice of timestamps representing the previous requests from the client,
// and returns whether the current request is allowed.
// If not allowed, it alos returns the time remaining until the next allowed request.
func AllowRequest(clientLogs *[]time.Time) (bool, time.Duration) {
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

	windowStartTime := time.Now().Add(time.Duration(-windowSize) * time.Second)
	firstRequestTime := (*clientLogs)[0]
	remainingTime := firstRequestTime.Sub(windowStartTime)
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

	curWindowStart := time.Now().Add(time.Duration(-windowSize) * time.Second)
	curTimestamp := (*clientLogs)[0]
	for len(*clientLogs) > 0 && curTimestamp.Before(curWindowStart) {
		*clientLogs = (*clientLogs)[1:]

		if len(*clientLogs) >= 1 {
			curTimestamp = (*clientLogs)[0]
		}
	}
}
