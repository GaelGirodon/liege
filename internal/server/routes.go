package server

import (
	"log/slog"
	"os"
	paths "path"
	"path/filepath"
	"sort"
	"strings"

	"gaelgirodon.fr/liege/internal/model"
)

// BuildRoutes loads stub response files from the given root directory and builds server routes.
func BuildRoutes(root string) (routes []*model.Route, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.Error("unable to access "+path, "error", err)
			return nil
		}
		if info.IsDir() {
			// Only serve regular files
			return nil
		}
		// Get the relative path to build the URL
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			slog.Error("unable to load "+path, "error", err)
			return nil
		}
		// Parse file name
		name, ext, route, err := parseFileName(info.Name())
		if err != nil {
			slog.Error("unable to load "+path, "error", err)
			return nil
		}
		// Build base URL
		baseUrl := strings.Trim(filepath.ToSlash(filepath.Dir(relPath)), "/.")
		// Load file and guess content type
		content, contentType, err := readFile(path)
		if err != nil {
			slog.Error("unable to read the file", "error", err)
			return nil
		}
		// 1st route: path without extension
		url := "/" + paths.Join(baseUrl, name)
		routes = append(routes, route.With(relPath, url, content, contentType))
		// 2nd route: full path?
		if len(ext) > 0 {
			routes = append(routes, route.With(relPath, url+ext, content, contentType))
		}
		// 3rd route: index file?
		if name == "index" {
			routes = append(routes, route.With(relPath, "/"+baseUrl, content, contentType))
		}
		return nil
	})
	// Sort routes by evaluation order
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Before(*routes[j])
	})
	return
}
