package controllers

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

type StaticController struct {
	FileServer http.Handler
	RootPath   string
}

func NewStaticController(rootPath string) *StaticController {
	return &StaticController{
		FileServer: http.FileServer(http.Dir(rootPath)),
		RootPath:   rootPath,
	}
}

func (ctrl *StaticController) ServeAssets(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(r.URL.Path)
	if !strings.HasPrefix(cleanPath, "/assets/") {
		http.Error(w, "Invalid path", http.StatusForbidden)
		return
	}

	fullPath := filepath.Join(ctrl.RootPath, cleanPath)
	if !strings.HasPrefix(filepath.Clean(fullPath), filepath.Clean(ctrl.RootPath)) {
		http.Error(w, "Invalid path", http.StatusForbidden)
		return
	}

	ctrl.FileServer.ServeHTTP(w, r)
}

func (ctrl *StaticController) ServeRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fullPath := filepath.Join(ctrl.RootPath, "index.html")
		if !strings.HasPrefix(filepath.Clean(fullPath), filepath.Clean(ctrl.RootPath)) {
			http.Error(w, "Invalid path", http.StatusForbidden)
			return
		}
		http.ServeFile(w, r, fullPath)
		return
	}
	ctrl.FileServer.ServeHTTP(w, r)
}
