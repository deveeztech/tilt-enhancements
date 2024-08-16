package await

import (
	"os"
	"strconv"
	"time"

	"github.com/deveeztech/tilt-enhancements/internal/netstat"
)

const (
	// DefaultPollingInterval is the default polling interval
	DefaultPollingInterval = 5
	// DefaultPortToListen is the default port to listen
	DefaultPortToListen = 40000

	// EnvPollingInterval overwrite the default polling interval in seconds
	EnvPollingInterval = "TILT_AWAIT_POLLING_INTERVAL"
	// EnvPortToListen overwrite the default port to listen
	EnvPortToListen = "TILT_AWAIT_PORT_TO_LISTEN"
)

type config struct {
	printf          func(string, ...interface{})
	pollingInterval time.Duration
	connType        string
	portToListen    int
}

func (c *config) log(fmt string, args ...interface{}) {
	if c.printf != nil {
		c.printf(fmt, args...)
	}
}

// Await until some process connecto to the desired port.
// Only works on Linux systems for now.
// netstat installed is not necessary.
func Start(opts ...Option) error {

	// Default configuration
	cfg := &config{
		pollingInterval: time.Duration(getEnv(EnvPollingInterval, DefaultPollingInterval)) * time.Second,
		connType:        netstat.TcpType,
		portToListen:    getEnv(EnvPortToListen, DefaultPortToListen),
	}

	for _, o := range opts {
		o.apply(cfg)
	}

	cfg.log("Starting await process in port %d ...", cfg.portToListen)
	for {

		d, err := netstat.Tcp()
		if err != nil {
			cfg.log("error getting tcp connections: %v", err)
			return err
		}

		for _, p := range d {
			if p.ForeignPort == cfg.portToListen && p.State == netstat.EstablishedState {
				cfg.log("new connection detected in port %d, exit from await", cfg.portToListen)
				return nil
			}
		}

		cfg.log("no conections to port %d, retrying in %s", cfg.portToListen, cfg.pollingInterval)
		time.Sleep(cfg.pollingInterval)

	}

}

func getEnv(key string, fallback int) int {

	if value, ok := os.LookupEnv(key); ok {
		result, err := strconv.Atoi(value)
		if err == nil {
			return result
		}
		// print an error message?
	}

	return fallback
}
