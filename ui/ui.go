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
	index, err := Asset("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t = template.Must(t.Parse(string(index)))
	indexTemplate = t
}

// Render will turn meetup event data into something to write out.
func Render(s data.Schedules) []byte {
	buf := &bytes.Buffer{}
	indexTemplate.Execute(buf, s)
	return buf.Bytes()
}
