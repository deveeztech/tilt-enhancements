package await

import (
	"time"

	"github.com/deveeztech/tilt-enhancements/internal/netstat"
)

type config struct {
	printf          func(string, ...interface{})
	pollingInterval time.Duration
	connType        string
	portToListen    int
}

// Await until some process connecto to the desired port.
// Only works on Linux systems for now.
// netstat installed is not necessary.
func Start(opts ...Option) error {

	// Default configuration
	cfg := &config{
		pollingInterval: 5 * time.Second,
		connType:        netstat.TCP,
		portToListen:    40000,
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
			if p.ForeignPort == cfg.portToListen && p.State == netstat.ESTABLISHED_STATE {
				cfg.log("new connection detected in port %d, exit from await", cfg.portToListen)
				return nil
			}
		}

		cfg.log("no conections to port %d, retrying in %s", cfg.portToListen, cfg.pollingInterval)
		time.Sleep(cfg.pollingInterval)

	}

}

func (c *config) log(fmt string, args ...interface{}) {
	if c.printf != nil {
		c.printf(fmt, args...)
	}
}
