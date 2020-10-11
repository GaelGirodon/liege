package server

import (
	"errors"
	"gaelgirodon.fr/liege/internal/console"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"os"
	paths "path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// A stub route configuration.
type Route struct {
	// Path to the loaded stub file.
	FilePath string `json:"file_path"`
	// URL path on which to serve this stub file.
	Path string `json:"path"`
	// Response status code.
	Code int `json:"code"`
	// Response body (stub file content).
	Content []byte `json:"-"`
	// Response content type.
	ContentType string `json:"content_type"`
}

// Load stub response files from the given root directory and build server routes.
func BuildRoutes(root string) (routes []*Route, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			console.Logger.Println("Error: unable to access " + path)
			return nil
		}
		if info.IsDir() {
			// Only serve regular files
			return nil
		}
		// Get the relative path to build the URL
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			console.Logger.Println("Error: unable to load " + path)
			return nil
		}
		// Parse file name
		name, ext, code := parseFileName(info.Name())
		// Build base URL
		baseUrl := strings.Trim(filepath.ToSlash(filepath.Dir(relPath)), "/.")
		// Load file and guess content type
		content, contentType, err := readFile(path)
		if err != nil {
			console.Logger.Println("Error: " + err.Error())
			return nil
		}
		// 1st route: path without extension
		url := "/" + paths.Join(baseUrl, name)
		routes = append(routes, &Route{relPath, url, code, content, contentType})
		// 2nd route: full path?
		if len(ext) > 0 {
			routes = append(routes, &Route{relPath, url + ext, code, content, contentType})
		}
		// 3rd route: index file?
		if name == "index" {
			routes = append(routes, &Route{relPath, "/" + baseUrl, code, content, contentType})
		}
		return nil
	})
	return
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
