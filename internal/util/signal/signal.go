package signal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Trap chan os.Signal

func NewTrap() Trap {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)

	return sigs
}

func (t Trap) Wait(ctx context.Context) {
	select {
	case sig := <-t:
		slog.Info("got signal", "signal", sig)
	case <-ctx.Done():
		slog.Info("context is done")
	}
	slog.Info("shutting down")
}
