package server

import (
	"log"
	"net/http"

	"github.com/milehighgophers/website/data"
	"github.com/milehighgophers/website/ui"
)

func Start(addr string) error {
	log.Printf("listening on %s", addr)

	mux := http.NewServeMux()
	mux.Handle("/", NewIndexHandler())
	return http.ListenAndServe(addr, mux)
}

type IndexHandler struct {
	store *data.Store
}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{
		store: &data.Store{},
	}
}

func (h *IndexHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	html := ui.Render(h.store.AllEvents())
	_, err := rw.Write(html)
	if err != nil {
		log.Print("error occured with /:", err)
	}
}
