package model

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Route is a stub route configuration.
type Route struct {
	// FilePath is the path to the loaded stub file.
	FilePath string `json:"file_path"`
	// Path is the URL path on which to serve this stub file.
	Path string `json:"path"`
	// Method is the required HTTP method.
	Method string `json:"method"`
	// QueryParams are the required query parameters.
	QueryParams []QueryParam `json:"query_params"`
	// Code is the response status code.
	Code int `json:"code"`
	// Content is the response body (stub file content).
	Content []byte `json:"-"`
	// ContentType is the response content type.
	ContentType string `json:"content_type"`
	// Latency is the simulated response latency (ms).
	Latency Latency `json:"latency"`
}

// NewRoute creates a new route structure With default values.
func NewRoute() Route {
	return Route{QueryParams: []QueryParam{}, Code: http.StatusOK, Latency: Latency{-1, -1}}
}

// With creates a new route structure with fields set.
func (r Route) With(filePath string, path string, content []byte, contentType string) *Route {
	route := r
	route.FilePath = filePath
	route.Path = path
	route.Content = content
	route.ContentType = contentType
	return &route
}

// Match checks the route eligibility against a given HTTP request.
func (r Route) Match(c echo.Context) bool {
	path := c.Request().URL.Path
	if r.Path != path || len(r.Method) > 0 && r.Method != c.Request().Method {
		return false
	}
	for _, qp := range r.QueryParams {
		if _, exists := c.QueryParams()[qp.Name]; !exists {
			return false
		}
		if len(qp.Value) > 0 && c.QueryParam(qp.Name) != qp.Value {
			return false
		}
	}
	return true
}

// Before reports whether the current route must be evaluated before the other one.
func (r Route) Before(r2 Route) bool {
	if r.Path != r2.Path { // Lexicographic order on path
		return r.Path < r2.Path
	}
	if r.Method != r2.Method { // Longer method before (empty/catch-all at the end)
		return len(r.Method) > len(r2.Method)
	}
	if len(r.QueryParams) != len(r2.QueryParams) { // Most query params before
		return len(r.QueryParams) > len(r2.QueryParams)
	}
	if len(r.QueryParams) > 0 { // Most-specific query params before
		return len(fmt.Sprintf("%v", r.QueryParams)) >
			len(fmt.Sprintf("%v", r2.QueryParams))
	}
	return r.FilePath < r2.FilePath // Lexicographic order on file path
}
