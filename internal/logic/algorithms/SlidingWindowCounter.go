package algorithms

import "time"

type SlidingWindowCounter struct {
	clientsWindows *map[string]map[string]*SlidingWindow
}

func (swc *SlidingWindowCounter) AllowRequest(clientKey, route string) (bool, time.Duration) {
	// If the route is not assigned in configs, we assume it is whitelisted
	if (*swc.clientsWindows)[route] == nil {
		return true, 0
	}

	swc.setupWindows(clientKey, route)

	clientWindow := (*swc.clientsWindows)[route][clientKey]

	clientWindow.checkWindow()

	weight := clientWindow.calculateWeight()

	if weight < float64(clientWindow.capacity) {
		clientWindow.currentRequests++
		return true, 0
	}

	remainingTime := time.Until(clientWindow.windowStart.Add(clientWindow.windowSize))
	return false, remainingTime
}

// setupWindows is a private function used as a helper function for the AllowRequest
// in Fixed Window Counter algorithm.
//
// It initializes or gets the client's window for the given route.
func (swc *SlidingWindowCounter) setupWindows(clientKey, route string) {
	_, exists := (*swc.clientsWindows)[route][clientKey]
	if !exists {
		(*swc.clientsWindows)[route][clientKey] = &SlidingWindow{
			capacity:        ConfigInstance.RoutesConfigs[route].Limit,
			currentRequests: 0,
			windowSize:      ConfigInstance.RoutesConfigs[route].Interval,
			windowStart:     time.Now(),
			previousCount:   0,
		}
	}
}

// checkWindow is a private function used to check if the sliding window
func (w *SlidingWindow) checkWindow() {
	if time.Since(w.windowStart) >= w.windowSize {
		w.previousCount = w.currentRequests
		w.currentRequests = 0
		w.windowStart = time.Now()
	}
}

// calculateWeight calculates the weighted request count based on the sliding window algorithm.
func (w *SlidingWindow) calculateWeight() float64 {
	elapsed := float64(time.Since(w.windowStart)) / float64(w.windowSize)
	threshold := (1-elapsed)*float64(w.previousCount) + float64(w.currentRequests) + 1
	return threshold
}
