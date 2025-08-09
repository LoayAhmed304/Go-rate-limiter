package algorithms

import (
	"time"
)

type TokenBucket struct {
	clientsBuckets *map[string]map[string]*Bucket
}

// AllowRequest determines whether a request may be allowed for a client,
// based on the Token Bucket algorithm.
//
// It takes a client key and a route to initialize or get the client's bucket for this route,
// and returns whether the request is allowed.
//
// If not allowed, it will also return the time remaining until the next allowed request.
func (tb *TokenBucket) AllowRequest(clientKey, route string) (bool, time.Duration) {
	// If the route is not assigned in configs, we assume it is whitelisted
	if (*tb.clientsBuckets)[route] == nil {
		return true, 0
	}

	tb.setupBuckets(clientKey, route)

	clientBucket := (*tb.clientsBuckets)[route][clientKey]

	clientBucket.refill()

	if clientBucket.curTokens > 0 {
		clientBucket.curTokens--
		return true, 0
	}

	remainingTime := clientBucket.refillInterval - time.Since(clientBucket.lastRefill)

	return false, remainingTime
}

func (tb *TokenBucket) Init(routes []string) {
	s := make(map[string]map[string]*Bucket)
	for _, route := range routes {
		s[route] = make(map[string]*Bucket)
	}
	tb.clientsBuckets = &s
}

// setupLogs is a private function used as a helper function for the AllowRequest
// in Token Bucket algorithm.
//
// It initializes or gets the client's buckets for the given route.
func (tb *TokenBucket) setupBuckets(clientKey, route string) {
	_, exists := (*tb.clientsBuckets)[route][clientKey]
	if !exists {
		(*tb.clientsBuckets)[route][clientKey] = &Bucket{
			refillInterval: ConfigInstance.RoutesConfigs[route].Interval / time.Duration(ConfigInstance.RoutesConfigs[route].Limit),
			capacity:       ConfigInstance.RoutesConfigs[route].Limit,
			curTokens:      ConfigInstance.RoutesConfigs[route].Limit,
			lastRefill:     time.Now(),
		}
	}
}

// refill is a method of the Bucket struct that refills the tokens based on the Bucket's refill interval.
//
// It calculates how many tokens to add based on the time since the last refill,
// by dividing the difference between now and the last refill time by the refill interval to determine how many tokens to add.
func (b *Bucket) refill() {
	curTime := time.Now()
	diff := curTime.Sub(b.lastRefill)
	tokensToAdd := int(diff / b.refillInterval)
	b.curTokens = min(b.capacity, b.curTokens+tokensToAdd)
	if tokensToAdd > 0 {
		b.lastRefill = curTime
	}
}
