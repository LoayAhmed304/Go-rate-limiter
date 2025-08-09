package algorithms

import (
	"time"
)

type TokenBucket struct {
	clientsBuckets *map[string]map[string]*Bucket
}

func (tb *TokenBucket) AllowRequest(clientIP, route string) (bool, time.Duration) {
	// If the route is not assigned in configs, we assume it is whitelisted
	if (*tb.clientsBuckets)[route] == nil {
		return true, 0
	}

	_, exists := (*tb.clientsBuckets)[route][clientIP]
	if !exists {
		(*tb.clientsBuckets)[route][clientIP] = &Bucket{
			refillInterval: ConfigInstance.RoutesConfigs[route].Interval / time.Duration(ConfigInstance.RoutesConfigs[route].Limit),
			capacity:       ConfigInstance.RoutesConfigs[route].Limit,
			curTokens:      ConfigInstance.RoutesConfigs[route].Limit,
			lastRefill:     time.Now(),
		}
	}
	clientBucket := (*tb.clientsBuckets)[route][clientIP]

	clientBucket.refill()

	if clientBucket.curTokens > 0 {
		clientBucket.curTokens -= 1
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

func (b *Bucket) refill() {
	curTime := time.Now()
	diff := curTime.Sub(b.lastRefill)
	tokensToAdd := int(diff / b.refillInterval)
	b.curTokens = min(b.capacity, b.curTokens+tokensToAdd)
	if tokensToAdd > 0 {
		b.lastRefill = curTime
	}
}
