package test

import (
	"errors"
	"fmt"
	"gaelgirodon.fr/liege/internal/model"
	"gaelgirodon.fr/liege/internal/server"
	"net"
	"testing"
	"time"
)

const (
	// root is the path to the stub files directory.
	root = "data"
	// port is the stub server HTTP port.
	port = 3000
)

// Test_e2e tests the application end-to-end
// (by sending requests to the server).
func Test_e2e(t *testing.T) {
	// Start the server asynchronously
	s := &server.StubServer{Config: model.Config{Root: root, Port: port}}
	go func() {
		_ = s.Start()
	}()
	// Wait for the server to be up
	err := errors.New("wait")
	for i := 0; err != nil && i < 10; i++ {
		time.Sleep(time.Second)
		_, err = net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Second)
	}

	// Test stub routes
	testStub(t)
	// Test management endpoints
	testManagementEndpoints(t)
}
