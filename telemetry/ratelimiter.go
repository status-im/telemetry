package telemetry

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const DEFAULT_CLEANUP_TIME = 1 * time.Hour

type Limiter struct {
	limiter  *rate.Limiter
	lastUsed time.Time
}

type RateLimiter struct {
	limiters map[string]*Limiter
	lock     *sync.RWMutex
	r        rate.Limit
	b        int
}

func NewRateLimiter(ctx context.Context, r rate.Limit, b int) *RateLimiter {
	return NewRateLimiterWithCleanup(ctx, r, b, DEFAULT_CLEANUP_TIME)
}

func NewRateLimiterWithCleanup(ctx context.Context, r rate.Limit, b int, cleanupTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*Limiter),
		lock:     &sync.RWMutex{},
		r:        r,
		b:        b,
	}

	go rl.cleanup(ctx, cleanupTime)

	return rl
}

func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.lock.Lock()

	limiter, ok := rl.limiters[ip]
	if !ok {
		rl.lock.Unlock()
		return rl.AddIP(ip)
	}

	limiter.lastUsed = time.Now()

	rl.lock.Unlock()
	return limiter.limiter
}

func (rl *RateLimiter) AddIP(ip string) *rate.Limiter {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	limiter := rate.NewLimiter(rl.r, rl.b)
	rl.limiters[ip] = &Limiter{limiter: limiter, lastUsed: time.Now()}

	return limiter
}

func (rl *RateLimiter) RemoveIP(ip string) {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	delete(rl.limiters, ip)
}

func (rl *RateLimiter) cleanup(ctx context.Context, cleanupEvery time.Duration) {
	t := time.NewTicker(cleanupEvery)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-t.C:
			for ip := range rl.limiters {
				limiter := rl.limiters[ip]
				if limiter.lastUsed.Add(2 * time.Second).Before(now) {
					rl.RemoveIP(ip)
				}
			}
		}
	}
}
