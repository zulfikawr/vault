package ui

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed dist/*
var distFS embed.FS

func Handler() http.Handler {
	stripped, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(stripped))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		
		// If the path doesn't have an extension, it's likely a client-side route
		// so serve index.html
		if !strings.Contains(path, ".") && path != "/" {
			path = "/"
		}

		// Manually set MIME types for embedded files because http.FS might not detect them correctly
		ext := filepath.Ext(path)
		if ext != "" {
			contentType := mime.TypeByExtension(ext)
			if contentType != "" {
				w.Header().Set("Content-Type", contentType)
			}
		}
		
		r.URL.Path = path
		fileServer.ServeHTTP(w, r)
	})
}