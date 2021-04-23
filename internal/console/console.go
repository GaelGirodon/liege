package console

import (
	"errors"
	"flag"
	"gaelgirodon.fr/liege/internal/model"
	"log"
	"math"
	"os"
	"strings"
)

const (
	// AppName is the application name.
	AppName = "liege"
	// Version is the application version number.
	Version = "0.3.0"
	// RootEnvVar is the name of the environment variable to set the root server directory.
	RootEnvVar = "LIEGE_ROOT"
	// PortEnvVar is the name of the environment variable to set the server port.
	PortEnvVar = "LIEGE_PORT"
	// LatencyEnvVar is the name of the environment variable to set the global latency.
	LatencyEnvVar = "LIEGE_LATENCY"
	// DefaultPort is the default HTTP server port number.
	DefaultPort = 3000
)

// Logger is the application global logger.
var Logger = log.New(os.Stdout, "", 0)

// ParseArgs parses and validates command-line args and environment vars.
func ParseArgs() (*model.Config, error) {
	// Print version number
	if len(os.Args) == 2 && os.Args[1] == "-v" {
		println("liege version " + Version)
		os.Exit(0)
	}
	// Parse args
	portFlag := flag.Uint("p", DefaultPort, "`port` to listen on")
	latencyFlag := flag.String("l", "0", "simulated response `latency` in ms")
	flag.Usage = func() {
		println("Usage:\n  " + AppName + " [flags] <root-dir>\n\nFlags:")
		flag.PrintDefaults()
	}
	args := strings.Join(os.Args, " ")
	// Default to environment variables
	if !strings.Contains(args, "-p") && len(os.Getenv(PortEnvVar)) > 0 {
		os.Args = append(os.Args, "-p="+os.Getenv(PortEnvVar))
	}
	if !strings.Contains(args, "-l") && len(os.Getenv(LatencyEnvVar)) > 0 {
		os.Args = append(os.Args, "-l="+os.Getenv(LatencyEnvVar))
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
	}
	if err := ValidateRootDirPath(root); err != nil {
		return nil, err
	}
	// Validate and parse latency
	latency, err := model.ParseLatency(*latencyFlag, "")
	if err != nil {
		return nil, errors.New("invalid latency value")
	}
	return &model.Config{Root: root, Port: uint16(*portFlag), Latency: latency}, nil
}

// ValidateRootDirPath checks that the given server root path points to a valid directory.
func ValidateRootDirPath(root string) error {
	if len(root) == 0 {
		return errors.New("path to the server root directory is required")
	} else if info, err := os.Stat(root); err != nil {
		return errors.New("root directory doesn't exist")
	} else if !info.IsDir() {
		return errors.New("root is not a directory")
	}
	return nil
}
