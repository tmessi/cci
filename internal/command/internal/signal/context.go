// Package signal is used to capture os signals.
package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// InitContext create a context that will cancel in the event of a SIGINT or SIGTERM.
func InitContext() (context.Context, context.CancelFunc) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigs
		cancel()
	}()
	return ctx, cancel
}
