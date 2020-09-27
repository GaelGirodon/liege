package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	// Application name.
	appName = "liege"
	// Application startup banner.
	banner = "_________ __   _________________\n" +
		"________ / /  /  _/ __/ ___/ __/\n" +
		"_______ / /___/ // _// (_ / _/\n" +
		"______ /____/___/___/\\___/___/"
	// Name of the environment variable to set the root server directory.
	rootEnvVar = "LIEGE_ROOT"
	// Name of the environment variable to set the server port.
	portEnvVar = "LIEGE_PORT"
	// Name of the environment variable to enable debug output.
	verboseEnvVar = "LIEGE_VERBOSE"
	// Default HTTP server port number.
	defaultPort = 3000
	// Request body with higher size won't be send back in a header.
	maxRequestBodySize = 4096
	// Request body header name in response.
	requestBodyHeader = "X-Request-Body"
)

// Program is running verbosely (display debug messages).
var verbose = false

func main() {
	// Parse command-line arguments
	root, port, v, err, code := parseArgs()
	verbose = v
	if err != nil {
		log("error", err.Error())
		os.Exit(code)
	}

	// Configure HTTP server
	log("", banner)
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	// Load stub files and register routes
	err = registerRoutes(root, e.Group(""))
	if err != nil {
		log("error", "unable to load server stub files ("+err.Error()+")")
		os.Exit(6)
	}

	// Start server
	log("info", "\nHTTP server started on port %d\n", port)
	e.Logger.Fatal(e.Start(":" + fmt.Sprint(port)))
}

// Parse, validate and return command-line arguments.
func parseArgs() (root string, port uint, verbose bool, err error, exitCode int) {
	// Parse args
	portFlag := flag.Uint("p", defaultPort, "port to listen on")
	verboseFlag := flag.Bool("v", false, "run verbosely")
	flag.Usage = func() {
		println("Usage:\n  " + appName + " [flags] <root-dir>\n\nFlags:")
		flag.PrintDefaults()
	}
	args := strings.Join(os.Args, " ")
	// Default to environment variable values
	if !strings.Contains(args, "-v") && len(os.Getenv(verboseEnvVar)) > 0 {
		os.Args = append(os.Args, "-v="+os.Getenv(verboseEnvVar))
	}
	if !strings.Contains(args, "-p") && len(os.Getenv(portEnvVar)) > 0 {
		os.Args = append(os.Args, "-p="+os.Getenv(portEnvVar))
	}
	flag.Parse()
	// Validate port
	if *portFlag < 80 || *portFlag > math.MaxUint16 {
		return "", 0, false, errors.New("invalid port number"), 2
	}
	// Validate root directory path
	root = os.Getenv(rootEnvVar)
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	} else if len(root) == 0 {
		return "", 0, false, errors.New("path to the server root directory is required"), 3
	}
	if info, err := os.Stat(root); err != nil {
		return "", 0, false, errors.New("root directory doesn't exist"), 4
	} else if !info.IsDir() {
		return "", 0, false, errors.New("root is not a directory"), 5
	}
	return root, *portFlag, *verboseFlag, nil, 0
}

// Load stub response files from the given root directory and register routes.
func registerRoutes(root string, server *echo.Group) error {
	log("debug", "\nServer root directory: "+root)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log("error", "unable to access "+path)
			return nil
		}
		if info.IsDir() {
			// Only serve regular files
			return nil
		}
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			log("error", "unable to load "+path)
			return nil
		}
		// Parse file name
		name, ext, code := parseFileName(info.Name())
		// Build URL
		url := "/" + strings.Trim(filepath.ToSlash(filepath.Dir(relPath)), "/")
		// Read file
		content, contentType, err := readFile(path)
		if err != nil {
			log("error", err.Error())
			return nil
		}
		// Register routes
		log("debug", relPath+":")
		handler := buildHandler(content, code, contentType)
		// Path without extension
		log("debug", "  %s => %d", url+"/"+name, code)
		server.Any(url+"/"+name, handler)
		// Full path
		if len(ext) > 0 {
			log("debug", "  %s => %d", url+"/"+name+ext, code)
			server.Any(url+"/"+name+ext, handler)
		}
		// Index file
		if name == "index" {
			log("debug", "  %s => %d", url, code)
			server.Any(url, handler)
		}
		return nil
	})
}

// Parse the file name and extract the name, extension and response status code.
func parseFileName(filename string) (name, ext string, code int) {
	ext = filepath.Ext(filename)
	name = strings.TrimSuffix(filename, ext)
	code = http.StatusOK
	if hasCode, _ := regexp.MatchString(".*__[1-5][0-9]{2}", name); hasCode {
		code, _ = strconv.Atoi(name[len(name)-3:])
		name = name[0 : len(name)-5]
	}
	return
}

// Read a file and return the contents and the content type (MIME type).
func readFile(path string) ([]byte, string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, "", errors.New("unable to read file")
	}
	contentType := ""
	if len(content) > 0 {
		contentType = http.DetectContentType(content)
	}
	if strings.Contains(contentType, "text/plain") {
		switch ext := filepath.Ext(path); ext {
		case ".json":
			contentType = strings.Replace(contentType, "text/plain", echo.MIMEApplicationJSON, 1)
		case ".xml":
			contentType = strings.Replace(contentType, "text/plain", echo.MIMETextXML, 1)
		case ".html":
			contentType = strings.Replace(contentType, "text/plain", echo.MIMETextHTML, 1)
		case ".js":
			contentType = strings.Replace(contentType, "text/plain", echo.MIMEApplicationJavaScript, 1)
		}
	}
	return content, contentType, nil
}

// Build route handler.
func buildHandler(content []byte, code int, contentType string) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Put original request body in an HTTP header
		var reqBody []byte
		if c.Request().Body != nil {
			reqBody, _ = ioutil.ReadAll(c.Request().Body)                 // Read
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset
			if len(reqBody) > 0 && len(reqBody) <= maxRequestBodySize {
				c.Response().Header().Set(requestBodyHeader, base64.StdEncoding.EncodeToString(reqBody))
			}
		}
		// Response
		if len(content) == 0 {
			return c.NoContent(code)
		}
		return c.Blob(code, contentType, content)
	}
}

// Log a message.
func log(level string, format string, a ...interface{}) {
	if level == "error" {
		println("Error: " + fmt.Sprintf(format, a...))
	} else if level != "debug" || verbose {
		fmt.Printf(format+"\n", a...)
	}
}
