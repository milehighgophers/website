package ui

import (
	"log"
	"net/http"
	"strings"
)

func NewAssetHandler() http.Handler {
	return &AssetHandler{}
}

type AssetHandler struct{}

func (*AssetHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// TODO: refactor this to pull data from assets more dynamically
	path := strings.TrimLeft(r.URL.Path, "/")
	switch path {
	case "assets/styles.css":
		sendAsset(path, "text/css", rw)
	case "assets/logo.png":
		sendAsset(path, "image/png", rw)
	case "assets/hero.jpg":
		sendAsset(path, "image/jpeg", rw)
	default:
		log.Printf("404 not found: %s", r.URL.Path)
		rw.WriteHeader(http.StatusNotFound)
	}
}

func sendAsset(name string, mimeType string, rw http.ResponseWriter) {
	file, err := Asset(name)
	if err != nil {
		log.Printf("failed to find asset: %s", file)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.Header().Add("content-type", mimeType)
	rw.Write(file)
}
