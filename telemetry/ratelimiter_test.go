package telemetry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func TestRateLimit(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	rl := NewRateLimiter(ctx, 1, 1, logger)

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

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	rl := NewRateLimiterWithCleanup(ctx, rate.Limit(1/time.Hour), 1, 100*time.Millisecond, logger)

	for i := 0; i < 300; i++ {
		ip := fmt.Sprintf("%d.%d.%d.%d", i, i, i, i)
		limiter := rl.GetLimiter(ip)
		require.True(t, limiter.Allow())
		require.False(t, limiter.Allow())
		time.Sleep(1 * time.Millisecond)
	}

	require.Equal(t, 300, rl.NumClients())

	time.Sleep(3 * time.Second)

	limiter2 := rl.GetLimiter("1.1.1.1")
	require.True(t, limiter2.Allow())

	require.Equal(t, 1, rl.NumClients())

}
