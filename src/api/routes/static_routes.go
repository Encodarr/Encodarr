package routes

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var allowedExtensions = map[string]bool{
	".html": true,
	".css":  true,
	".js":   true,
	".png":  true,
	".jpg":  true,
	".ico":  true,
}

func HandleStatic(dir string) (http.HandlerFunc, http.HandlerFunc) {
	cleanDir := filepath.Clean(dir)
	absDir, err := filepath.Abs(cleanDir)
	if err != nil {
		panic(err)
	}

	// Assets handler
	assetsHandler := func(w http.ResponseWriter, r *http.Request) {
		requestPath := filepath.Clean(r.URL.Path)

		// Check for directory traversal
		if strings.Contains(requestPath, "..") {
			log.Printf("Rejected path traversal attempt: %s", requestPath)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Validate file extension
		ext := strings.ToLower(filepath.Ext(requestPath))
		if !allowedExtensions[ext] {
			log.Printf("Rejected invalid extension: %s", ext)
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}

		fullPath := filepath.Join(absDir, requestPath)

		// Ensure path is within allowed directory
		if !strings.HasPrefix(fullPath, absDir) {
			log.Printf("Rejected path outside root: %s", fullPath)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Resolve any symbolic links
		realPath, err := filepath.EvalSymlinks(fullPath)
		if err != nil || !strings.HasPrefix(realPath, absDir) {
			log.Printf("Rejected symlink path: %s", fullPath)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Check if file exists
		if _, err := os.Stat(realPath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, realPath)
	}

	// Root handler for SPA
	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(absDir, "index.html"))
	}

	return assetsHandler, rootHandler
}
