package test

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// testStub tests routes created from test stub files.
func testStub(t *testing.T) {
	jsonHeaders := map[string]string{"Content-Type": "application/json; charset=utf-8"}
	tests := []struct {
		name        string
		method      string
		path        string
		wantStatus  int
		wantBody    string
		wantHeaders map[string]string
		wantLatency time.Duration
	}{
		// items/index.json
		{"e2e/index/get/200/1", http.MethodGet, "/items", http.StatusOK, "[]", jsonHeaders, 0},
		{"e2e/index/get/200/2", http.MethodGet, "/items/index", http.StatusOK, "[]", jsonHeaders, 0},
		{"e2e/index/get/200/3", http.MethodGet, "/items/index.json", http.StatusOK, "[]", jsonHeaders, 0},
		{"e2e/index/post/200", http.MethodPost, "/items", http.StatusOK, "[]", jsonHeaders, 0},
		{"e2e/index/put/200", http.MethodPut, "/items", http.StatusOK, "[]", jsonHeaders, 0},
		{"e2e/index/patch/200", http.MethodPatch, "/items", http.StatusOK, "[]", jsonHeaders, 0},
		// items/index__qs.json
		{"e2e/query/get/200/1", http.MethodGet, "/items?s=1", http.StatusOK, "[{\"id\": 1}]", jsonHeaders, 0},
		{"e2e/query/get/200/2", http.MethodGet, "/items/index?s=2", http.StatusOK, "[{\"id\": 1}]", jsonHeaders, 0},
		{"e2e/query/get/200/3", http.MethodGet, "/items/index.json?s=3", http.StatusOK, "[{\"id\": 1}]", jsonHeaders, 0},
		{"e2e/query/delete/200", http.MethodDelete, "/items/index.json?s=put", http.StatusOK, "[{\"id\": 1}]", jsonHeaders, 0},
		// items/1__GET.json
		{"e2e/single/get/200/1", http.MethodGet, "/items/1", http.StatusOK, "{}", jsonHeaders, 0},
		{"e2e/single/get/200/2", http.MethodGet, "/items/1.json", http.StatusOK, "{}", jsonHeaders, 0},
		{"e2e/single/post/404", http.MethodPost, "/items/1", http.StatusNotFound, "", nil, 0},
		// admin/index__403_l50
		{"e2e/forbidden/get/403/1", http.MethodGet, "/admin", http.StatusForbidden, "", nil, 50 * time.Millisecond},
		{"e2e/forbidden/get/403/2", http.MethodGet, "/admin/index", http.StatusForbidden, "", nil, 50 * time.Millisecond},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			start := time.Now()
			// Prepare and send the request to the stub server
			url := fmt.Sprintf("http://localhost:%d%s", port, test.path)
			var body io.Reader
			if test.method != http.MethodGet && test.wantStatus == http.StatusOK {
				// Pass a request body to be able to test the X-Request-Body response header
				body = strings.NewReader(test.wantBody)
			} else {
				body = http.NoBody
			}
			req, _ := http.NewRequest(test.method, url, body)
			res, _ := http.DefaultClient.Do(req)
			// Check the response status
			if res.StatusCode != test.wantStatus {
				t.Errorf("want status = %d, got %d", test.wantStatus, res.StatusCode)
			}
			// Check the latency
			lat := time.Since(start)
			if test.wantLatency > 0 && (lat < test.wantLatency || lat > test.wantLatency*2) {
				t.Errorf("want latency ~= %d, got %d", test.wantLatency, lat)
			}
			// Check the response body
			if len(test.wantBody) > 0 {
				if body, err := io.ReadAll(res.Body); err != nil {
					t.Fatalf("Unexpected error reading body: %s", err.Error())
				} else if strings.TrimSpace(string(body)) != test.wantBody {
					t.Errorf("want body = %q, got %q", test.wantBody, string(body))
				}
			}
			_ = res.Body.Close()
			// Check response headers
			for key, value := range test.wantHeaders {
				header := res.Header.Get(key)
				if res.Header.Get(key) != value {
					t.Errorf("want header %q = %q, got %q", key, value, header)
				}
			}
			// Check the X-Request-Body response header
			if test.method != http.MethodGet && test.wantStatus == http.StatusOK {
				reqBodyBase64 := res.Header.Get("X-Request-Body")
				if len(reqBodyBase64) == 0 {
					t.Errorf("want X-Request-Body to be set")
				} else if reqBody, err := base64.StdEncoding.DecodeString(reqBodyBase64); err != nil {
					t.Errorf("want X-Request-Body to be correctly base64 encoded")
				} else if string(reqBody) != test.wantBody {
					t.Errorf("want X-Request-Body = %q, got %q", test.wantBody, string(reqBody))
				}
			}
		})
	}
}
