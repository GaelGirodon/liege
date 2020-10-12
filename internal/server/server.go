package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"gaelgirodon.fr/liege/internal/console"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"net/http"
)

const (
	// Request body with higher size won't be send back in a header.
	maxRequestBodySize = 4096
	// Request body header name in response.
	requestBodyHeader = "X-Request-Body"
)

// An HTTP server for stub files.
type StubServer struct {
	// Root server directory.
	Root string
	// Server port number.
	Port uint16
	// Stub routes.
	routes []*Route
}

// Start the stub server.
func (s *StubServer) Start() error {
	// Setup HTTP server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	// Load stub files and build routes
	routes, err := BuildRoutes(s.Root)
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
	console.Logger.Printf("\nHTTP server started on port %d\n\n", s.Port)
	e.Logger.Fatal(e.Start(":" + fmt.Sprint(s.Port)))
	return nil
}

//
// Handlers
//

// Get the stub server configuration.
func (s *StubServer) getConfigHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, console.Args{Root: s.Root})
}

// Update the stub server configuration.
func (s *StubServer) updateConfigHandler(c echo.Context) error {
	config := new(console.Args)
	if err := c.Bind(config); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	} else if err := console.ValidateRootDirPath(config.Root); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	s.Root = config.Root
	return c.NoContent(http.StatusNoContent)
}

// Reload stub files and re-build routes.
func (s *StubServer) refreshHandler(c echo.Context) error {
	routes, err := BuildRoutes(s.Root)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			"unable to load stub files and build routes: "+err.Error())
	}
	s.routes = routes
	return c.NoContent(http.StatusNoContent)
}

// Return current registered routes.
func (s *StubServer) routesHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.routes)
}

// Handle stub requests using the registered routes.
func (s *StubServer) stubsHandler(c echo.Context) error {
	url := c.Request().URL.Path
	for _, route := range s.routes {
		if route.Path != url {
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
		if len(route.Content) == 0 {
			return c.NoContent(route.Code)
		}
		return c.Blob(route.Code, route.ContentType, route.Content)
	}
	return c.NoContent(http.StatusNotFound)
}
