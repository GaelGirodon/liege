package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gaelgirodon.fr/liege/internal/console"
	"gaelgirodon.fr/liege/internal/model"
)

const (
	// Request body with size higher than maxRequestBodySize
	// won't be sent back in a header.
	maxRequestBodySize = 4096
	// requestBodyHeader is the request body header name in the response.
	requestBodyHeader = "X-Request-Body"
)

// StubServer is an HTTP server for stub files.
type StubServer struct {
	// Config is the application configuration.
	Config model.Config
	// routes it the stub routes list.
	routes []*model.Route
}

// Start starts the stub server.
func (s *StubServer) Start() error {
	// Load stub files and build routes
	routes, err := BuildRoutes(s.Config.Root)
	if err != nil {
		return err
	}
	s.routes = routes
	// Register first-level routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /_liege/config", s.getConfigHandler)
	mux.HandleFunc("PUT /_liege/config", s.updateConfigHandler)
	mux.HandleFunc("POST /_liege/refresh", s.refreshHandler)
	mux.HandleFunc("GET /_liege/routes", s.routesHandler)
	mux.HandleFunc("/", s.stubsHandler)
	// Start
	h := &http.Server{
		Addr:    s.Config.Address(),
		Handler: s.middleware(mux),
	}
	if s.Config.HasTLS() {
		slog.Info("HTTPS server started on port " + strconv.Itoa(int(s.Config.Port)))
		err = h.ListenAndServeTLS(s.Config.Cert, s.Config.Key)
	} else {
		slog.Info("HTTP server started on port " + strconv.Itoa(int(s.Config.Port)))
		err = h.ListenAndServe()
	}
	return err
}

//
// Middleware
//

// middleware handles common features on all handlers.
func (s *StubServer) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Log request
		slog.Info("REQUEST", "method", req.Method, "uri", req.RequestURI,
			"host", req.Host, "user_agent", req.UserAgent(), "remote_addr", req.RemoteAddr)
		// Remove trailing slash
		if len(req.URL.Path) > 1 && strings.HasSuffix(req.URL.Path, "/") {
			req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
			req.RequestURI = req.URL.RequestURI()
		}
		// Next
		next.ServeHTTP(w, req)
	})
}

//
// Handlers
//

// getConfigHandler gets the stub server configuration.
func (s *StubServer) getConfigHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(s.Config); err != nil {
		slog.Error("unable to write config data to the response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// updateConfigHandler updates the stub server configuration.
func (s *StubServer) updateConfigHandler(w http.ResponseWriter, req *http.Request) {
	config := new(model.Config)
	if err := json.NewDecoder(req.Body).Decode(config); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
	} else if err := console.ValidateRootDirPath(config.Root); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if !config.Latency.IsValid() {
		http.Error(w, "invalid latency value", http.StatusBadRequest)
	} else {
		s.Config.Root = config.Root
		s.Config.Latency = config.Latency
		w.WriteHeader(http.StatusNoContent)
	}
}

// refreshHandler reloads stub files and re-builds routes.
func (s *StubServer) refreshHandler(w http.ResponseWriter, _ *http.Request) {
	if routes, err := BuildRoutes(s.Config.Root); err != nil {
		http.Error(w, "unable to load stub files and build routes: "+err.Error(),
			http.StatusBadRequest)
	} else {
		s.routes = routes
		w.WriteHeader(http.StatusNoContent)
	}
}

// routesHandler returns current registered routes.
func (s *StubServer) routesHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(s.routes); err != nil {
		slog.Error("unable to write routes data to the response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// stubsHandler handles stub requests using the registered routes.
func (s *StubServer) stubsHandler(w http.ResponseWriter, req *http.Request) {
	for _, route := range s.routes {
		if !route.Match(req) {
			continue
		}
		var reqBody []byte
		if req.Body != nil && req.ContentLength > 0 {
			reqBody, _ = io.ReadAll(req.Body)                           // Read request body
			req.Body = io.NopCloser(bytes.NewBuffer(reqBody))           // Reset request body
			if len(reqBody) > 0 && len(reqBody) <= maxRequestBodySize { // Set as a response header
				w.Header().Set(requestBodyHeader, base64.StdEncoding.EncodeToString(reqBody))
			}
		}
		latency := route.Latency.Compute(s.Config.Latency)
		if latency > 0 {
			time.Sleep(latency)
		}
		if len(route.Content) == 0 {
			w.WriteHeader(route.Code)
		} else {
			w.Header().Set("Content-Type", route.ContentType)
			w.WriteHeader(route.Code)
			if _, err := w.Write(route.Content); err != nil {
				slog.Error("unable to write stub route content to the response")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
