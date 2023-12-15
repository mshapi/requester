package service

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"sync"

	"golang.org/x/time/rate"
	"requester/internal/model"
)

//go:generate mockery --name=HTTPClient --inpackage --case=snake
type HTTPClient interface {
	Post(ctx context.Context, url, body string)
}

func New(client HTTPClient) *Requester {
	return &Requester{client: client}
}

type Requester struct {
	client HTTPClient
}

func (r *Requester) Run(ctx context.Context, req *model.RequestData) error {
	if err := validateURL(req.URL); err != nil {
		return err
	}

	if req.Amount <= 0 || req.PerSecond == 0 {
		return nil
	}

	limiter := rate.NewLimiter(rate.Limit(req.PerSecond), req.PerSecond)
	wg := &sync.WaitGroup{}
	wg.Add(req.Amount)

	for i := 0; i < req.Amount; i++ {
		if err := limiter.Wait(ctx); err != nil {
			return err
		}

		go r.doReq(ctx, req.URL, i, wg)
	}

	wg.Wait()

	return nil
}

// makeBody return string like as `{ "iteration": 0 }`
func (r *Requester) makeBody(i int) string {
	return `{ "iteration": ` + strconv.Itoa(i) + ` }`
}

func (r *Requester) doReq(ctx context.Context, url string, i int, wg *sync.WaitGroup) {
	r.client.Post(ctx, url, r.makeBody(i))
	wg.Done()
}

func validateURL(link string) error {
	u, err := url.Parse(link)
	if u.Scheme == "" {
		return errors.New("empty scheme")
	}
	if u.Host == "" {
		return errors.New("empty host")
	}
	return err
}