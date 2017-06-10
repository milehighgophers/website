package ui

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type AssetHandler struct {
	mTypes map[string]string
}

func NewAssetHandler() http.Handler {
	return &AssetHandler{
		mTypes: map[string]string{
			"css":  "text/css",
			"png":  "image/png",
			"jpg":  "image/jpeg",
			"jpeg": "image/jpeg",
		},
	}
}

func (h *AssetHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	data, err := Asset(path)
	if err != nil {
		log.Printf("404 not found: %s", r.URL.Path)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Header().Add("content-type", h.mimeType(path))
	rw.Write(data)
}

func (h *AssetHandler) mimeType(path string) string {
	ext := filepath.Ext(path)
	return h.mTypes[ext]
}
