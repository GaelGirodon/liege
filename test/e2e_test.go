package test

import (
	"gaelgirodon.fr/liege/internal/model"
	"gaelgirodon.fr/liege/internal/server"
	"testing"
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

	// Test stub routes
	testStub(t)
	// Test management endpoints
	testManagementEndpoints(t)
}
