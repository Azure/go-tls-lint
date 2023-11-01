package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
)

// Provider provides the CLI settings.
type Provider interface {
	// Context - creates a context for the command to run in.
	Context() (context.Context, context.CancelFunc)

	// Logger - creates a logger for the command to use.
	Logger() *slog.Logger

	// Printf - prints a message to the user.
	Printf(format string, args ...interface{})
}

type cliProvider struct {
	Verbose bool
}


var _ Provider = (*cliProvider)(nil)

func (*cliProvider) Context() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt)
}

func (p *cliProvider) Logger() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelInfo,
	}
	if p.Verbose {
		opts.Level = slog.LevelDebug
	}

	return slog.New(slog.NewTextHandler(os.Stderr, opts))
}

func (*cliProvider) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
