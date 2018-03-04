package main

import (
	"time"

	"k8s.io/client-go/util/workqueue"
)

// FixedItemIntervalRateLimiter limits items to a fixed-rate interval
type FixedItemIntervalRateLimiter struct {
	interval time.Duration
}

var _ workqueue.RateLimiter = &FixedItemIntervalRateLimiter{}

func NewFixedItemIntervalRateLimiter(interval time.Duration) workqueue.RateLimiter {
	return &FixedItemIntervalRateLimiter{
		interval: interval,
	}
}

func (r *FixedItemIntervalRateLimiter) When(item interface{}) time.Duration {
	return r.interval
}

func (r *FixedItemIntervalRateLimiter) NumRequeues(item interface{}) int {
	return 1
}

func (r *FixedItemIntervalRateLimiter) Forget(item interface{}) {
}

// NewDefaultCAUSRateLimitter creates a rate limitter which limits overall (as per the
// default controller rate limiter), as well as per the resync interval
func NewDefaultCAUSRateLimiter(interval time.Duration) workqueue.RateLimiter {
	return NewFixedItemIntervalRateLimiter(interval)
}
