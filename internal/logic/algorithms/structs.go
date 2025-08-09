package algorithms

import "time"

type Algorithm interface {
	AllowRequest(string, string) (bool, time.Duration)
	Init([]string)
}

type RouteConfig struct {
	Limit    int           `json:"limit"`
	Interval time.Duration `json:"interval"`
}

type Config struct {
	Algorithm     Algorithm              `json:"algorithm"`
	RoutesConfigs map[string]RouteConfig `json:"routes"`
}

type rawConfig struct {
	Algorithm string           `json:"algorithm"`
	Routes    []rawRouteConfig `json:"routes"`
}

type rawRouteConfig struct {
	Route    string `json:"route"`
	Limit    int    `json:"limit"`
	Interval string `json:"interval"`
}

type Bucket struct {
	refillInterval time.Duration
	capacity       int
	curTokens      int
	lastRefill     time.Time
}

type Window struct {
	capacity        int
	currentRequests int
	windowSize      time.Duration
	windowStart     time.Time
}
