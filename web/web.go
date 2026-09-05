package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SpaHandler() http.HandlerFunc {
	distDir := "dist"
	fileServer := http.FileServer(http.Dir(distDir))

	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		fullPath := filepath.Join(distDir, filepath.Clean(path))

		// Check if file exists on disk
		info, err := os.Stat(fullPath)
		if err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fallback to index.html for SPA routes
		http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
	}
}
