package main

import (
	"gaelgirodon.fr/liege/internal/console"
	"gaelgirodon.fr/liege/internal/server"
)

// main is the application entrypoint.
func main() {
	// Parse command-line arguments
	cfg, err := console.ParseArgs()
	if err != nil {
		console.Logger.Fatalln("Error: " + err.Error())
	}

	// Print start-up banner
	console.Logger.Println("_________ __   _________________\n" +
		"________ / /  /  _/ __/ ___/ __/\n" +
		"_______ / /___/ // _// (_ / _/\n" +
		"______ /____/___/___/\\___/___/")

	// Setup and start the HTTP server
	s := &server.StubServer{Config: *cfg}
	if err = s.Start(); err != nil {
		console.Logger.Fatalln("Error: " + err.Error())
	}
}
