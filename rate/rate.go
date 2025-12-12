package rate

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

	r.advance()
	if r.tokens == 1 {
		r.tokens--
		return true
	}
	return false
}

func (r *RateLimiter) Take() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.advance()
	time.Sleep(r.durationFromTokens(1 - r.tokens)) // 0 returns immediately
	r.tokens--
}

func (r *RateLimiter) advance() {
	now := time.Now()
	delta := now.Sub(r.last).Seconds() * r.limit
	r.tokens = min(r.tokens+delta, 1) // burst of 1
	r.last = now
}

func (r *RateLimiter) durationFromTokens(tokens float64) time.Duration {
	return time.Duration((tokens / r.limit) * float64(time.Second))
}
