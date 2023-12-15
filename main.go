package main

import (
	"context"
	"os/signal"
	"syscall"

	"requester/internal/command"
	resty_client "requester/internal/pkg/resty-client"
	"requester/internal/service"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	requester := service.New(resty_client.New())

	command.New(
		ctx,
		requester,
	).Execute()
}
