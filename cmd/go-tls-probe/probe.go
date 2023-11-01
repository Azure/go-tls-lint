package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

type Probe struct {
	Verbose bool `short:"v" help:"Enable verbose logging."`

	Addr string `arg:"" help:"host[:port] to probe."`
}

func (p *Probe) Run(provider Provider) error {
	ctx, cancel := provider.Context()
	defer cancel()

	logger := provider.Logger().With("addr", p.Addr)
	startAt := time.Now()

	dialer := &tls.Dialer{
		NetDialer: new(net.Dialer),
		Config: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	logger.DebugContext(ctx, "dialing")
	c, err := dialer.DialContext(ctx, "tcp", p.Addr) // handshake has happened here
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", p.Addr, err)
	}
	logger.DebugContext(ctx, "connected")

	conn := c.(*tls.Conn)
	defer conn.Close()

	connState := conn.ConnectionState()
	provider.Printf("Connected to %s (duration=%s)\n", p.Addr, time.Since(startAt))
	provider.Printf("TLS version: %s\n", tls.VersionName(connState.Version))
	provider.Printf("Cipher suite: %s\n", tls.CipherSuiteName(connState.CipherSuite))

	return nil
}
