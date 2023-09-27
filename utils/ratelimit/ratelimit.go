package util

import (
	"context"
	"time"

	"github.com/ragpanda/go-toolkit/log"
	"go.uber.org/ratelimit"
)

type RateLimit struct {
	limit    int
	duration time.Duration
	name     string
}

// NewRateLimit create a ratelimit
// wrapper raw uber pkg for metric and log
func NewRateLimit(name string, limit int, duration time.Duration) *RateLimit {
	return &RateLimit{
		limit:    limit,
		duration: duration,
		name:     name,
	}
}
func (self *RateLimit) Take(ctx context.Context) {
	r := ratelimit.New(self.limit, ratelimit.Per(self.duration))
	start := time.Now()
	_ = r.Take()

	cost := time.Now().Sub(start)
	if cost > 5*time.Second {
		log.Warn(ctx, "wait too long, %s %v", self.name, cost.String())
	}

	return
}
