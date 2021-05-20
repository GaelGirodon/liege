package server

import (
	"bytes"
	"encoding/base64"
	"gaelgirodon.fr/liege/internal/console"
	"gaelgirodon.fr/liege/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"net/http"
	"time"
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
	// Setup HTTP server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	// Load stub files and build routes
	routes, err := BuildRoutes(s.Config.Root)
	if err != nil {
		return err
	}
	s.routes = routes
	// Register first-level routes
	e.GET("/_liege/config", s.getConfigHandler)
	e.PUT("/_liege/config", s.updateConfigHandler)
	e.POST("/_liege/refresh", s.refreshHandler)
	e.GET("/_liege/routes", s.routesHandler)
	e.Any("/*", s.stubsHandler)
	// Start
	if s.Config.HasTLS() {
		console.Logger.Printf("\nHTTPS server started on port %d\n\n", s.Config.Port)
		err = e.StartTLS(s.Config.Address(), s.Config.Cert, s.Config.Key)
	} else {
		console.Logger.Printf("\nHTTP server started on port %d\n\n", s.Config.Port)
		err = e.Start(s.Config.Address())
	}
	e.Logger.Fatal(err)
	return nil
}

//
// Handlers
//

// getConfigHandler gets the stub server configuration.
func (s *StubServer) getConfigHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.Config)
}

// updateConfigHandler updates the stub server configuration.
func (s *StubServer) updateConfigHandler(c echo.Context) error {
	config := new(model.Config)
	if err := c.Bind(config); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	} else if err := console.ValidateRootDirPath(config.Root); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if !config.Latency.IsValid() {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid latency value")
	}
	s.Config.Root = config.Root
	s.Config.Latency = config.Latency
	return c.NoContent(http.StatusNoContent)
}

// refreshHandler reloads stub files and re-builds routes.
func (s *StubServer) refreshHandler(c echo.Context) error {
	if routes, err := BuildRoutes(s.Config.Root); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			"unable to load stub files and build routes: "+err.Error())
	} else {
		s.routes = routes
	}
	return c.NoContent(http.StatusNoContent)
}

// routesHandler returns current registered routes.
func (s *StubServer) routesHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.routes)
}

// stubsHandler handles stub requests using the registered routes.
func (s *StubServer) stubsHandler(c echo.Context) error {
	for _, route := range s.routes {
		if !route.Match(c) {
			continue
		}
		var reqBody []byte
		if c.Request().Body != nil && c.Request().ContentLength > 0 {
			reqBody, _ = ioutil.ReadAll(c.Request().Body)                 // Read request body
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset request body
			if len(reqBody) > 0 && len(reqBody) <= maxRequestBodySize {   // Set as a response header
				c.Response().Header().Set(requestBodyHeader, base64.StdEncoding.EncodeToString(reqBody))
			}
		}
		latency := route.Latency.Compute(s.Config.Latency)
		if latency > 0 {
			time.Sleep(latency)
		}
		if len(route.Content) == 0 {
			return c.NoContent(route.Code)
		}
		return c.Blob(route.Code, route.ContentType, route.Content)
	}
	return c.NoContent(http.StatusNotFound)
}
