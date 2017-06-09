package ui

import (
	"bytes"
	"html/template"

	"github.com/milehighgophers/website/data"
)

var indexTemplate = template.Must(
	template.ParseFiles(
		"ui/templates/index.html",
		"ui/templates/events.html",
	),
)

// Render will turn meetup event data into something to write out.
func Render(s *data.MeetupSchedule) []byte {
	buf := &bytes.Buffer{}
	indexTemplate.Execute(buf, s)
	return buf.Bytes()
}
