package await

import "time"

// An Option alters the behavior of Init.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (of optionFunc) apply(cfg *config) { of(cfg) }

// Logger uses the supplied printf implementation for log output. By default,
// Set doesn't log anything.
func Logger(printf func(string, ...interface{})) Option {
	return optionFunc(func(cfg *config) {
		cfg.printf = printf
	})
}

// ConnectionType sets the type of connection to await.
func ConnectionType(connType string) Option {
	return optionFunc(func(cfg *config) {
		cfg.connType = connType
	})
}

// PollingInterval sets the interval between polling attempts.
func PollingInterval(d time.Duration) Option {
	return optionFunc(func(cfg *config) {
		cfg.pollingInterval = d
	})
}

// PortToListen sets the port to await connections on.
func PortToListen(port int) Option {
	return optionFunc(func(cfg *config) {
		cfg.portToListen = port
	})
}
