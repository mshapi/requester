package resty_client

import (
	"context"

	"github.com/go-resty/resty/v2"
)

func New() *client {
	return &client{}
}

type client struct{}

func (c *client) Post(ctx context.Context, url, body string) {
	_, _ = resty.New().
		R().
		SetContext(ctx).
		SetBody(body).
		Post(url)
}
