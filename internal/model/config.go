package model

// Config is the application configuration.
type Config struct {
	// Root is the path to the root server directory.
	Root string `json:"root"`
	// Port is the HTTP server port number.
	Port uint16 `json:"-"`
	// Latency is the simulated response latency value.
	Latency Latency `json:"latency"`
}
