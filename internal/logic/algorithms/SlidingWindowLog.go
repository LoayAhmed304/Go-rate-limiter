package algorithms

import (
	"time"
)

var windowSize time.Duration
var maxRequests int

type SlidingWindowLog struct {
	ClientsLogs *map[string]map[string]*[]time.Time
}

// AllowRequest determines whether a request may be allowed for a client,
// based on the Sliding Window Log algorithm.
//
// It takes a slice of timestamps representing the previous requests from the client,
// and returns whether the current request is allowed.
// If not allowed, it alos returns the time remaining until the next allowed request.
func (swl *SlidingWindowLog) AllowRequest(clientIP, route string) (bool, time.Duration) {
	// route is not assigned in configs, so we assume it is whitelisted
	if (*swl.ClientsLogs)[route] == nil {
		return true, 0
	}

	windowSize = ConfigInstance.RoutesConfigs[route].Interval
	maxRequests = ConfigInstance.RoutesConfigs[route].Limit

	_, exists := (*swl.ClientsLogs)[route][clientIP]
	if !exists {
		s := make([]time.Time, 0, maxRequests)
		(*swl.ClientsLogs)[route][clientIP] = &s
	}

	clientLogs := (*swl.ClientsLogs)[route][clientIP]

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

func (swl *SlidingWindowLog) Init(routes []string) {
	s := make(map[string]map[string]*[]time.Time)

	for _, route := range routes {
		s[route] = make(map[string]*[]time.Time)
	}

	swl.ClientsLogs = &s
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
