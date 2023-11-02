package basic

import "crypto/tls"

func emptyFunc() {}

func emptyConfig() {
	_ = &tls.Config{}
	_ = tls.Config{}
}

func configSettings() {
	a := tls.Config{
		Renegotiation: tls.RenegotiateFreelyAsClient, // want `WARNING: Unexpected TLS config settings "Renegotiation"`
		MinVersion:    tls.VersionTLS12,              // want `ERROR: Go 1.18 and onward has good defaults for MinVersion, no need to set it`
	}
	a.Renegotiation = tls.RenegotiateFreelyAsClient // want `WARNING: Unexpected TLS config settings "Renegotiation"`
	a.MinVersion = tls.VersionTLS12                 // want `ERROR: Go 1.18 and onward has good defaults for MinVersion, no need to set it`

	b := &tls.Config{
		Renegotiation: tls.RenegotiateFreelyAsClient, // want `WARNING: Unexpected TLS config settings "Renegotiation"`
		MinVersion:    tls.VersionTLS12,              // want `ERROR: Go 1.18 and onward has good defaults for MinVersion, no need to set it`
	}
	b.Renegotiation = tls.RenegotiateFreelyAsClient // want `WARNING: Unexpected TLS config settings "Renegotiation"`
	b.MinVersion = tls.VersionTLS12                 // want `ERROR: Go 1.18 and onward has good defaults for MinVersion, no need to set it`
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
