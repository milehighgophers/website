package server

import (
	"log"
	"net/http"

	"github.com/milehighgophers/website/data"
	"github.com/milehighgophers/website/ui"
)

type Store interface {
	AllEvents() map[string][]data.Event
}

func Start(addr string, s Store) error {
	log.Printf("listening on %s", addr)

	mux := http.NewServeMux()
	mux.Handle("/", NewIndexHandler(s))
	return http.ListenAndServe(addr, mux)
}

type IndexHandler struct {
	store Store
}

func NewIndexHandler(s Store) *IndexHandler {
	return &IndexHandler{
		store: s,
	}
}

func (h *IndexHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	html := ui.Render(h.store.AllEvents())
	_, err := rw.Write(html)
	if err != nil {
		log.Print("error occured with /:", err)
	}
}
