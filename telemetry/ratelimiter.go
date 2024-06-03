package telemetry

import (
	"context"
	"sync"
	"time"

	shrinkingmap "github.com/go-auxiliaries/shrinking-map/pkg/shrinking-map"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

const DEFAULT_CLEANUP_TIME = 1 * time.Hour

type Limiter struct {
	limiter  *rate.Limiter
	lastUsed time.Time
}

type RateLimiter struct {
	limiters *shrinkingmap.Map[string, *Limiter] //map[string]*Limiter
	lock     *sync.RWMutex
	r        rate.Limit
	b        int
	logger   *zap.Logger
}

func NewRateLimiter(ctx context.Context, r rate.Limit, b int, logger *zap.Logger) *RateLimiter {
	return NewRateLimiterWithCleanup(ctx, r, b, DEFAULT_CLEANUP_TIME, logger)
}

func NewRateLimiterWithCleanup(ctx context.Context, r rate.Limit, b int, cleanupTime time.Duration, logger *zap.Logger) *RateLimiter {
	rl := &RateLimiter{
		limiters: shrinkingmap.New[string, *Limiter](200),
		lock:     &sync.RWMutex{},
		r:        r,
		b:        b,
		logger:   logger,
	}

	go rl.cleanup(ctx, cleanupTime)

	return rl
}

func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.lock.Lock()

	limiter, ok := rl.limiters.Get2(ip)
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
	rl.limiters.Set(ip, &Limiter{limiter: limiter, lastUsed: time.Now()})

	return limiter
}

func (rl *RateLimiter) RemoveIP(ip string) {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	rl.limiters.Delete(ip)
}

func (rl *RateLimiter) NumClients() int {
	return len(rl.limiters.Values())
}

func (rl *RateLimiter) cleanup(ctx context.Context, cleanupEvery time.Duration) {
	t := time.NewTicker(cleanupEvery)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-t.C:
			numCleaned := 0
			for ip, limiter := range rl.limiters.Values() {
				if limiter.lastUsed.Add(2 * time.Second).Before(now) {
					rl.RemoveIP(ip)
					numCleaned++
				}
			}
			rl.logger.Debug("cleanup", zap.Int("removed", numCleaned))
		}
	}
}
