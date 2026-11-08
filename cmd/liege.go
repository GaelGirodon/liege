package main

import (
	"fmt"
	"log/slog"
	"os"

	"gaelgirodon.fr/liege/internal/console"
	"gaelgirodon.fr/liege/internal/server"
)

// main is the application entrypoint.
func main() {
	// Parse command-line arguments
	cfg, err := console.ParseArgs()
	if err != nil {
		slog.Error("unable to parse command-line arguments", "error", err)
		os.Exit(1)
	}

	// Print start-up banner
	fmt.Printf("_________ __   _________________\n" +
		"________ / /  /  _/ __/ ___/ __/\n" +
		"_______ / /___/ // _// (_ / _/\n" +
		"______ /____/___/___/\\___/___/\n\n")

	// Setup and start the HTTP server
	s := &server.StubServer{Config: *cfg}
	if err = s.Start(); err != nil {
		slog.Error("unable to start the stub server", "error", err)
		os.Exit(1)
	}
}
