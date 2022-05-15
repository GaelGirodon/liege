package test

import (
	"encoding/json"
	"fmt"
	"gaelgirodon.fr/liege/internal/model"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

// testManagementEndpoints tests /_liege/* endpoints.
func testManagementEndpoints(t *testing.T) {
	// GET /_liege/config => get and check configuration
	t.Run("e2e/mngmt/config/get", func(t *testing.T) {
		checkConfigEndpoint(t, model.Config{Root: root, Latency: model.Latency{Min: 0, Max: 0}})
	})

	// PUT /_liege/config => try to update the configuration with an invalid latency
	t.Run("e2e/mngmt/config/put/400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/_liege/config", port),
			strings.NewReader(fmt.Sprintf(`{"root":"%s","latency":{"min":3,"max":2}}`, root)))
		req.Header.Set("Content-Type", "application/json")
		res, _ := http.DefaultClient.Do(req)
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("want status = %d, got %d", http.StatusBadRequest, res.StatusCode)
		}
	})

	// PUT /_liege/config => update and check configuration
	t.Run("e2e/mngmt/config/put/204", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/_liege/config", port),
			strings.NewReader(fmt.Sprintf(`{"root":"%s","latency":{"min":1,"max":2}}`, root)))
		req.Header.Set("Content-Type", "application/json")
		res, _ := http.DefaultClient.Do(req)
		if res.StatusCode != http.StatusNoContent {
			t.Errorf("want status = %d, got %d", http.StatusNoContent, res.StatusCode)
		}
		checkConfigEndpoint(t, model.Config{Root: root, Latency: model.Latency{Min: 1, Max: 2}})
	})

	// GET /_liege/routes => get and check routes
	t.Run("e2e/mngmt/routes/get", func(t *testing.T) {
		checkRoutesEndpoint(t, 10)
	})

	// POST /_liege/refresh => modify & reload stub files and check routes
	t.Run("e2e/mngmt/refresh/post", func(t *testing.T) {
		_ = os.WriteFile("data/test", []byte(""), 0666)
		res, _ := http.Post(fmt.Sprintf("http://localhost:%d/_liege/refresh", port), "", http.NoBody)
		if res.StatusCode != http.StatusNoContent {
			t.Errorf("want status = %d, got %v", http.StatusNoContent, res.StatusCode)
		}
		checkRoutesEndpoint(t, 11)
		_ = os.Remove("data/test")
	})
}

// checkConfigEndpoint requests the /_liege/config endpoint
// and compares the response body with the given configuration.
func checkConfigEndpoint(t *testing.T, wantConfig model.Config) {
	res, _ := http.Get(fmt.Sprintf("http://localhost:%d/_liege/config", port))
	// Check the response status
	if res.StatusCode != http.StatusOK {
		t.Errorf("want status = %d, got %d", http.StatusOK, res.StatusCode)
	}
	// Check the response body
	body, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	wantBody := fmt.Sprintf(`{"root":"%s","latency":{"min":%d,"max":%d}}`,
		wantConfig.Root, wantConfig.Latency.Min, wantConfig.Latency.Max)
	if strings.TrimSpace(string(body)) != wantBody {
		t.Errorf("want body = %s, got %s", wantBody, body)
	}
}

// checkRoutesEndpoint requests the /_liege/routes endpoint
// and checks the routes count in the response body.
func checkRoutesEndpoint(t *testing.T, wantRoutes int) {
	res, _ := http.Get(fmt.Sprintf("http://localhost:%d/_liege/routes", port))
	// Check the response status
	if res.StatusCode != http.StatusOK {
		t.Errorf("want status = %d, got %d", http.StatusOK, res.StatusCode)
	}
	// Check the response body
	body, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	var routes []model.Route
	_ = json.Unmarshal(body, &routes)
	if len(routes) != wantRoutes {
		t.Errorf("want %d routes, got %d", wantRoutes, len(routes))
	}
}
