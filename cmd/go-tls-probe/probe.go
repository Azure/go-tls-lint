package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
)

type Probe struct {
	Verbose bool `short:"v" help:"Enable verbose logging."`

	Addr string `arg:"" help:"host[:port] to probe."`
}

func (p *Probe) Run() error {
	ctx := context.Background()

	dialer := &tls.Dialer{
		NetDialer: new(net.Dialer),
		Config: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}
	c, err := dialer.DialContext(ctx, "tcp", p.Addr) // handshake has happened here
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", p.Addr, err)
	}
	conn := c.(*tls.Conn)
	defer conn.Close()

	connState := conn.ConnectionState()
	fmt.Println(tls.VersionName(connState.Version), tls.CipherSuiteName(connState.CipherSuite))

	return nil
}
