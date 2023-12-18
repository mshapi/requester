package limiter

import (
	"context"
	"math"
	"testing"
	"time"
)

func TestLimiter_Wait(t *testing.T) {
	tests := []struct {
		name       string
		ctx        func() context.Context
		limit      int64
		period     time.Duration
		iterations int64
		wantErr    bool
	}{
		{
			name:       "ok",
			ctx:        context.Background,
			limit:      4,
			period:     time.Second,
			iterations: 18,
			wantErr:    false,
		},
		{
			name: "canceled context",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			limit:      4,
			period:     time.Second,
			iterations: 18,
			wantErr:    true,
		},
		{
			name: "canceled context in process",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				return ctx
			},
			limit:      10,
			period:     time.Second,
			iterations: 120,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.limit, tt.period)

			ctx := tt.ctx()

			start := time.Now()

			var err error

			for i := int64(0); i < tt.iterations; i++ {
				err = l.Wait(ctx)
				if err != nil {
					break
				}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("expected want error: %t, but given error: %#v", tt.wantErr, err)
				return
			}

			sub := time.Now().Sub(start)
			need := time.Duration(math.Ceil(float64(tt.iterations) / float64(tt.limit)))

			if sub < need {
				t.Errorf(
					"time elapsed: %s, must be min: %s\n",
					sub,
					need,
				)
			}
		})
	}
}
