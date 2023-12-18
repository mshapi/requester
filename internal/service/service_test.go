package service

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
	"requester/internal/model"
)

func TestNew(t *testing.T) {
	client := NewMockHTTPClient(t)

	r := New(client)

	require.Equal(t, reflect.ValueOf(client).Pointer(), reflect.ValueOf(r.client).Pointer())
}

func TestRequester_Run(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		client  func(t *testing.T, req *model.RequestData) HTTPClient
		ctx     context.Context
		req     *model.RequestData
		wantErr bool
	}{
		{
			name: "invalid url",
			client: func(t *testing.T, req *model.RequestData) HTTPClient {
				return nil
			},
			req: &model.RequestData{
				URL:       "://asd",
				Amount:    10,
				PerSecond: 10,
			},
			wantErr: true,
		},
		{
			name: "invalid url scheme",
			client: func(t *testing.T, req *model.RequestData) HTTPClient {
				return nil
			},
			req: &model.RequestData{
				URL:       "asd",
				Amount:    10,
				PerSecond: 10,
			},
			wantErr: true,
		},
		{
			name: "invalid url host",
			client: func(t *testing.T, req *model.RequestData) HTTPClient {
				return nil
			},
			req: &model.RequestData{
				URL:       "asd://",
				Amount:    10,
				PerSecond: 10,
			},
			wantErr: true,
		},
		{
			name: "null amount",
			client: func(t *testing.T, req *model.RequestData) HTTPClient {
				return nil
			},
			req: &model.RequestData{
				URL:       "http://asd.e",
				Amount:    10,
				PerSecond: 0,
			},
			wantErr: false,
		},
		{
			name: "null per second",
			client: func(t *testing.T, req *model.RequestData) HTTPClient {
				return nil
			},
			req: &model.RequestData{
				URL:       "http://asd.e",
				Amount:    0,
				PerSecond: 10,
			},
			wantErr: false,
		},
		{
			name: "context canceled",
			client: func(t *testing.T, req *model.RequestData) HTTPClient {
				return nil
			},
			ctx: (func() context.Context {
				tmp, cancel := context.WithCancel(ctx)
				cancel()
				return tmp
			})(),
			req: &model.RequestData{
				URL:       "http://asd.e",
				Amount:    5,
				PerSecond: 2,
			},
			wantErr: true,
		},
		{
			name: "check rps",
			client: func(t *testing.T, req *model.RequestData) HTTPClient {
				c := NewMockHTTPClient(t)

				limiter := rate.NewLimiter(rate.Limit(req.PerSecond), req.PerSecond)

				rMx := sync.Mutex{}
				requests := make(map[int]struct{})
				for i := 0; i < req.Amount; i++ {
					requests[i] = struct{}{}
				}

				type reqBody struct {
					Iteration int `json:"iteration"`
				}

				failed := atomic.Bool{}

				t.Cleanup(func() {
					if failed.Load() {
						return
					}

					rMx.Lock()
					defer rMx.Unlock()
					if len(requests) > 0 {
						panic("not all iterations have been completed")
					}
				})

				c.On("Post", ctx, req.URL, mock.Anything).Run(func(args mock.Arguments) {
					defer func() {
						if err := recover(); err != nil {
							failed.Store(true)
							panic(err)
						}
					}()

					if !limiter.Allow() {
						panic("the permissible RPS value has been exceeded")
					}

					tmp := &reqBody{}
					if err := json.NewDecoder(strings.NewReader(args.String(2))).Decode(tmp); err != nil {
						panic(err)
					}

					rMx.Lock()
					defer rMx.Unlock()

					delete(requests, tmp.Iteration)
				}).Times(req.Amount)

				return c
			},
			req: &model.RequestData{
				URL:       "http://asd.e",
				Amount:    16,
				PerSecond: 3,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Requester{
				client: tt.client(t, tt.req),
			}
			ctx := ctx
			if tt.ctx != nil {
				ctx = tt.ctx
			}
			if err := r.Run(ctx, tt.req); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
