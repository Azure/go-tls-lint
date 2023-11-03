package basic

import (
	"crypto/tls"
	"time"
)

func emptyFunc() {}

func emptyConfig() {
	_ = &tls.Config{}
	_ = tls.Config{}
}

func configSettings() {
	a := tls.Config{
		Renegotiation: tls.RenegotiateFreelyAsClient, // want `WARNING: Renegotiation is not supported in TLS 1.3`
		MinVersion:    tls.VersionTLS12,              // want `ERROR: Go 1.18 and onward has good defaults for MinVersion, no need to set it`
		Time:          time.Now,                      // should not raise issue
	}
	a.Renegotiation = tls.RenegotiateFreelyAsClient // want `WARNING: Renegotiation is not supported in TLS 1.3`
	a.MinVersion = tls.VersionTLS12                 // want `ERROR: Go 1.18 and onward has good defaults for MinVersion, no need to set it`
	a.Time = time.Now                               // should not raise issue

	b := &tls.Config{
		Renegotiation: tls.RenegotiateFreelyAsClient, // want `WARNING: Renegotiation is not supported in TLS 1.3`
		MinVersion:    tls.VersionTLS12,              // want `ERROR: Go 1.18 and onward has good defaults for MinVersion, no need to set it`
		Time:          time.Now,                      // should not raise issue
	}
	b.Renegotiation = tls.RenegotiateFreelyAsClient // want `WARNING: Renegotiation is not supported in TLS 1.3`
	b.MinVersion = tls.VersionTLS12                 // want `ERROR: Go 1.18 and onward has good defaults for MinVersion, no need to set it`
	b.Time = time.Now                               // should not raise issue
}

func configSettingsWithComment() {
	a := tls.Config{
		// go-tls-lint:ignore - ignore this line
		Renegotiation: tls.RenegotiateFreelyAsClient,
		// go-tls-lint:ignore - this line should not be ignored
		MinVersion: tls.VersionTLS12, // want `ERROR: Go 1.18 and onward has good defaults for MinVersion, no need to set it`
	}
	// go-tls-lint:ignore - ignore this line
	a.Renegotiation = tls.RenegotiateFreelyAsClient
	// go-tls-lint:ignore - this line should not be ignored
	a.MinVersion = tls.VersionTLS12 // want `ERROR: Go 1.18 and onward has good defaults for MinVersion, no need to set it`

	b := &tls.Config{
		// go-tls-lint:ignore - ignore this line
		Renegotiation: tls.RenegotiateFreelyAsClient,
		// go-tls-lint:ignore - this line should not be ignored
		MinVersion: tls.VersionTLS12, // want `ERROR: Go 1.18 and onward has good defaults for MinVersion, no need to set it`
	}
	// go-tls-lint:ignore - ignore this line
	a.Renegotiation = tls.RenegotiateFreelyAsClient
	// go-tls-lint:ignore - this line should not be ignored
	b.Renegotiation = tls.RenegotiateFreelyAsClient
	b.MinVersion = tls.VersionTLS12 // want `ERROR: Go 1.18 and onward has good defaults for MinVersion, no need to set it`
}
