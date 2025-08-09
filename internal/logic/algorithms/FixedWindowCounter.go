package algorithms

import "time"

type FixedWindowCounter struct {
	clientsWindows *map[string]map[string]*Window
}

// AllowRequest determines whether a request may be allowed for a client,
// based on the Fixed Window Counter algorithm.
//
// It takes a client key and a route to initialize or get the client's window for this route,
// and returns whether the request is allowed.
//
// If not allowed, it will also return the time remaining until the next allowed request.
func (fwc *FixedWindowCounter) AllowRequest(clientKey, route string) (bool, time.Duration) {
	// If the route is not assigned in configs, we assume it is whitelisted
	if (*fwc.clientsWindows)[route] == nil {
		return true, 0
	}

	fwc.setupWindows(clientKey, route)

	clientWindow := (*fwc.clientsWindows)[route][clientKey]

	clientWindow.checkWindow()

	if clientWindow.currentRequests < clientWindow.capacity {
		clientWindow.currentRequests++
		return true, 0
	}

	remainingTime := time.Until(clientWindow.windowStart.Add(clientWindow.windowSize))
	return false, remainingTime
}

func (fwc *FixedWindowCounter) Init(routes []string) {
	s := make(map[string]map[string]*Window)
	for _, route := range routes {
		s[route] = make(map[string]*Window)
	}
	fwc.clientsWindows = &s
}

// setupWindows is a private function used as a helper function for the AllowRequest
// in Fixed Window Counter algorithm.
//
// It initializes or gets the client's window for the given route.
func (fwc *FixedWindowCounter) setupWindows(clientKey, route string) {
	_, exists := (*fwc.clientsWindows)[route][clientKey]
	if !exists {
		(*fwc.clientsWindows)[route][clientKey] = &Window{
			capacity:        ConfigInstance.RoutesConfigs[route].Limit,
			currentRequests: 0,
			windowSize:      ConfigInstance.RoutesConfigs[route].Interval,
			windowStart:     time.Now(),
		}
	}
}

// checkWindow checks if the current window has expired and resets the request count if it has.
func (w *Window) checkWindow() {
	if time.Since(w.windowStart) >= w.windowSize {
		w.currentRequests = 0
		w.windowStart = time.Now()
	}
}
