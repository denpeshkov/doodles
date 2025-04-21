package main

import (
	"sync"
	"time"
)

type RateLimiter struct {
	limit float64

	mu     sync.Mutex
	tokens float64
	last   time.Time
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{limit: float64(limit)}
}

func (r *RateLimiter) CanTake() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	r.tokens = r.advance(now.Sub(r.last))
	r.last = now

	if r.tokens == 1 {
		r.tokens--
		return true
	}
	return false
}

func (r *RateLimiter) Take() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	r.tokens = r.advance(now.Sub(r.last))
	r.last = now

	time.Sleep(r.durationFromTokens(1 - r.tokens)) // 0 returns immediately
	r.tokens--
}

func (r *RateLimiter) advance(dur time.Duration) (newTokens float64) {
	delta := dur.Seconds() * r.limit
	return min(r.tokens+delta, 1) // burst of 1
}

func (r *RateLimiter) durationFromTokens(tokens float64) time.Duration {
	return time.Duration((tokens / r.limit) * float64(time.Second))
}
