package model

import "fmt"

// Config is the application configuration.
type Config struct {
	// Root is the path to the root server directory.
	Root string `json:"root"`
	// Port is the HTTP server port number.
	Port uint16 `json:"-"`
	// Cert is the path to the TLS certificate PEM file.
	Cert string `json:"-"`
	// Key is the path to the TLS private key PEM file.
	Key string `json:"-"`
	// Latency is the simulated response latency value.
	Latency Latency `json:"latency"`
}

// Address returns the HTTP server address.
func (c *Config) Address() string {
	return ":" + fmt.Sprint(c.Port)
}

// HasTLS indicates whether TLS configuration is provided or not.
func (c *Config) HasTLS() bool {
	return len(c.Cert) > 0 && len(c.Key) > 0
}
