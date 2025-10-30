package internal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func NewContextWithSignal(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		ch := make(chan os.Signal, 1)

		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(ch)

		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
			return
		}
	}()
	return ctx
}
