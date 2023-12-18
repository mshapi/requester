package limiter

import (
	"context"
	"sync/atomic"
	"time"
)

func New(limit int64, period time.Duration) *Limiter {
	return &Limiter{
		limit:  limit,
		period: period,
	}
}

type Limiter struct {
	limit  int64
	period time.Duration

	current atomic.Int64

	start atomic.Int64 // unix nano
}

func (l *Limiter) init() {
	l.start.CompareAndSwap(0, time.Now().UnixNano())
	l.current.CompareAndSwap(0, l.limit)
}

func (l *Limiter) reset() *Limiter {
	l.start.Store(time.Now().UnixNano())
	l.current.Store(l.limit)
	return l
}

func (l *Limiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	l.init()

	if l.current.Add(-1) > 0 {
		return nil
	}

	select {
	case <-time.After(time.Unix(0, l.start.Load()).Add(l.period).Sub(time.Now())):
		l.reset()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
