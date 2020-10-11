package console

import (
	"errors"
	"flag"
	"log"
	"math"
	"os"
	"strings"
)

const (
	// Application name.
	AppName = "liege"
	// Name of the environment variable to set the root server directory.
	RootEnvVar = "LIEGE_ROOT"
	// Name of the environment variable to set the server port.
	PortEnvVar = "LIEGE_PORT"
	// Default HTTP server port number.
	DefaultPort = 3000
)

// Application global logger.
var Logger = log.New(os.Stdout, "", 0)

// Command-line arguments.
type Args struct {
	// Root server directory.
	Root string
	// Server port.
	Port uint16
}

// Parse, validate and return command-line arguments.
func ParseArgs() (*Args, error) {
	// Parse args
	portFlag := flag.Uint("p", DefaultPort, "`port` to listen on")
	flag.Usage = func() {
		println("Usage:\n  " + AppName + " [flags] <root-dir>\n\nFlags:")
		flag.PrintDefaults()
	}
	args := strings.Join(os.Args, " ")
	// Default to environment variable values
	if !strings.Contains(args, "-p") && len(os.Getenv(PortEnvVar)) > 0 {
		os.Args = append(os.Args, "-p="+os.Getenv(PortEnvVar))
	}
	flag.Parse()
	// Validate port
	if *portFlag < 80 || *portFlag > math.MaxUint16 {
		return nil, errors.New("invalid port number")
	}
	// Validate root directory path
	root := os.Getenv(RootEnvVar)
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	} else if len(root) == 0 {
		return nil, errors.New("path to the server root directory is required")
	}
	if info, err := os.Stat(root); err != nil {
		return nil, errors.New("root directory doesn't exist")
	} else if !info.IsDir() {
		return nil, errors.New("root is not a directory")
	}
	return &Args{root, uint16(*portFlag)}, nil
}
