package main

import (
	"gaelgirodon.fr/liege/internal/console"
	"gaelgirodon.fr/liege/internal/server"
)

func main() {
	// Parse command-line arguments
	args, err := console.ParseArgs()
	if err != nil {
		console.Logger.Fatalln("Error: " + err.Error())
	}

	// Start-up banner
	console.Logger.Println("_________ __   _________________\n" +
		"________ / /  /  _/ __/ ___/ __/\n" +
		"_______ / /___/ // _// (_ / _/\n" +
		"______ /____/___/___/\\___/___/")

	// Setup and start the HTTP server
	s := &server.StubServer{Root: args.Root, Port: args.Port}
	if err = s.Start(); err != nil {
		console.Logger.Fatalln("Error: " + err.Error())
	}
}
