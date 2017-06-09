package ui

import (
	"log"
	"net/http"

	"github.com/milehighgophers/website/ui/assets"
)

func NewAssetHandler() http.Handler {
	return &AssetHandler{}
}

type AssetHandler struct{}

func (*AssetHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/assets/styles.css":
		sendAsset("styles.css", "text/css", rw)
	case "/assets/logo.png":
		sendAsset("logo.png", "image/png", rw)
	case "/assets/hero.jpg":
		sendAsset("hero.jpg", "image/jpeg", rw)
	default:
		log.Printf("404 not found: %s", r.URL.Path)
		rw.WriteHeader(http.StatusNotFound)
	}
}

func sendAsset(name string, mimeType string, rw http.ResponseWriter) {
	file, err := assets.Asset(name)
	if err != nil {
		log.Printf("failed to find asset: %s", file)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.Header().Add("content-type", mimeType)
	rw.Write(file)
}
