package server

import (
	"errors"
	"gaelgirodon.fr/liege/internal/model"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	// optsPrefix is the prefix of options in a file name.
	optsPrefix = "__"
	// optsSeparator is the separator between options in a file name.
	optsSeparator = "_"
)

var (
	// methodOptPattern is the pattern to match the HTTP method option.
	methodOptPattern = regexp.MustCompile("(?i)^(GET|HEAD|POST|PUT|PATCH|DELETE|CONNECT|OPTIONS|TRACE)$")
	// queryOptPattern is the pattern to match a request query parameter option.
	queryOptPattern = regexp.MustCompile("^q([a-z0-9-]+)(?:=([a-z0-9-]+))?$")
	// codeOptPattern is the pattern to match the custom HTTP response status code option.
	codeOptPattern = regexp.MustCompile("^([1-5][0-9]{2})$")
)

// parseFileName parses the file name and extract the name, extension and options.
func parseFileName(filename string) (name, ext string, route model.Route, err error) {
	ext = filepath.Ext(filename)
	name = strings.TrimSuffix(filename, ext)
	route = model.NewRoute()
	optsPrefixIndex := strings.Index(name, optsPrefix)
	if optsPrefixIndex == -1 {
		return // No options
	}
	// Parse options
	opts := strings.Split(name[optsPrefixIndex+len(optsPrefix):], optsSeparator)
	name = name[:optsPrefixIndex]
	for _, opt := range opts {
		if len(opt) == 0 {
			continue
		} else if match := methodOptPattern.FindStringSubmatch(opt); len(match) == 2 {
			route.Method = match[1]
		} else if match := queryOptPattern.FindStringSubmatch(opt); len(match) == 3 {
			route.QueryParams = append(route.QueryParams, model.QueryParam{Name: match[1], Value: match[2]})
		} else if match := codeOptPattern.FindStringSubmatch(opt); len(match) == 2 {
			route.Code, _ = strconv.Atoi(match[1])
		} else if latency, parsingErr := model.ParseLatency(opt, "l"); parsingErr == nil {
			route.Latency = latency
		} else {
			err = errors.New("unknown or invalid option '" + opt + "'")
			return
		}
	}
	return
}

// readFile reads a file and returns the contents and the content type (MIME type).
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
