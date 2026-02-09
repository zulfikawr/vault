package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed dist/*
var distFS embed.FS

func Handler() http.Handler {
	// Get the sub-filesystem from the dist directory
	stripped, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(stripped))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the path doesn't have an extension, it's likely a client-side route
		// so serve index.html
		if !strings.Contains(r.URL.Path, ".") && r.URL.Path != "/" {
			r.URL.Path = "/"
		}
		
		fileServer.ServeHTTP(w, r)
	})
}