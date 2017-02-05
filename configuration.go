package astiredis

import "flag"

// Flags
var (
	Addr   = flag.String("redis-addr", "", "the redis addr")
	Prefix = flag.String("redis-prefix", "", "the redis prefix")
)

// Configuration represents the configuration of the proxy
type Configuration struct {
	Addr   string `toml:"addr"`
	Prefix string `toml:"prefix"`
}

// FlagConfig generates a Configuration based on flags
func FlagConfig() Configuration {
	return Configuration{
		Addr:   *Addr,
		Prefix: *Prefix,
	}
}
