package telemetry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestRateLimit(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()
	rl := NewRateLimiter(ctx, 1, 1)

	ip1 := "1.1.1.1"

	limiter := rl.GetLimiter(ip1)
	require.True(t, limiter.Allow())

	limiter = rl.GetLimiter(ip1)
	require.False(t, limiter.Allow())

	time.Sleep(1 * time.Second)
	limiter = rl.GetLimiter(ip1)
	require.True(t, limiter.Allow())

	ip2 := "2.2.2.2:8080"
	limiter = rl.GetLimiter(ip2)
	require.True(t, limiter.Allow())

	limiter = rl.GetLimiter(ip2)
	require.False(t, limiter.Allow())
}

func TestRateLimitCleanup(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()
	rl := NewRateLimiterWithCleanup(ctx, rate.Limit(1/time.Hour), 1, 100*time.Millisecond)

	ip1 := "1.1.1.1"

	limiter := rl.GetLimiter(ip1)
	require.True(t, limiter.Allow())
	require.False(t, limiter.Allow())

	time.Sleep(3 * time.Second)

	limiter2 := rl.GetLimiter(ip1)
	require.True(t, limiter2.Allow())
}
