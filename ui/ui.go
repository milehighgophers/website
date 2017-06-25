package ui

import (
	"bytes"
	"html/template"
	"log"

	"github.com/milehighgophers/website/data"
)

//go:generate go-bindata -pkg ui assets templates

var indexTemplate *template.Template

func init() {
	t := template.New("index")
	events, err := Asset("templates/events.html")
	if err != nil {
		log.Fatal(err)
	}
	index, err := Asset("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t = template.Must(t.Parse(string(index)))
	t = template.Must(t.Parse(string(events)))
	indexTemplate = t
}

// Render will turn meetup event data into something to write out.
func Render(s *data.MeetupSchedule) []byte {
	buf := &bytes.Buffer{}
	indexTemplate.Execute(buf, s)
	return buf.Bytes()
}
