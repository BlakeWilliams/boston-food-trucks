package static

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/blakewilliams/medium/pkg/router"
)

type Config struct {
	// Path to the directory containing the static files.
	FileRoot string
	// Prefix for routes requesting static resources.
	// e.g. "/static" would serve files from "/static/<file>"
	PathPrefix string
}

// New returns a middleware handler that serves static files from the given
// directory based on Config.
func Middleware(config Config) router.Middleware {
	fileRoot, err := filepath.Abs(config.FileRoot)
	if err != nil {
		panic(err)
	}

	return func(c router.Action, next router.MiddlewareFunc) {
		if !strings.HasPrefix(c.Request().URL.Path, config.PathPrefix) {
			next(c)
			return
		}

		relativePath := strings.TrimPrefix(c.Request().URL.Path, config.PathPrefix)
		relativePath = strings.TrimPrefix(relativePath, "/")
		absPath, err := filepath.Abs(filepath.Join(fileRoot, relativePath))

		if err != nil {
			c.Response().WriteHeader(http.StatusNotFound)
			c.Response().Write([]byte("404 not found"))
			return
		}

		if !strings.HasPrefix(absPath, fileRoot) {
			fmt.Println("wtf")
			c.Response().WriteHeader(http.StatusNotFound)
			c.Response().Write([]byte("404 not found"))
			return
		}

		http.ServeFile(c.Response(), c.Request(), absPath)
	}
}
