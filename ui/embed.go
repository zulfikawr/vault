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

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// If the path doesn't have an extension, it's likely a client-side route
		// so serve index.html
		if !strings.Contains(path, ".") && path != "/" {
			path = "/"
		}

		// Manually set MIME types for embedded files
		ext := filepath.Ext(path)
		if ext != "" {
			contentType := mime.TypeByExtension(ext)
			if contentType != "" {
				w.Header().Set("Content-Type", contentType)
			}
		}

		cleanPath := strings.TrimPrefix(path, "/")
		
		// Handle root path
		if cleanPath == "" {
			cleanPath = "index.html"
		}
		
		file, err := stripped.Open(cleanPath)
		if err != nil {
			// Try to serve index.html if path is a directory or not found
			indexPath := cleanPath
			if !strings.HasSuffix(indexPath, "/index.html") {
				if indexPath != "" && indexPath != "index.html" {
					indexPath = cleanPath + "/index.html"
				} else {
					indexPath = "index.html"
				}
			}
			
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			data, err := fs.ReadFile(stripped, indexPath)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Write(data)
			return
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if info.IsDir() {
			// Try to serve index.html from the directory
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			indexPath := cleanPath + "/index.html"
			data, err := fs.ReadFile(stripped, indexPath)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Write(data)
		} else {
			data, err := fs.ReadFile(stripped, cleanPath)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Write(data)
		}
	})
}
