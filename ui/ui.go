package ui

import (
	"bytes"
	"html/template"

	"github.com/milehighgophers/website/data"
)

const (
	indexTemplateStr = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Mile High Gopher Events</title>
	</head>
	<body>
		{{range .}}
			<div>{{.Name}} -- {{.Time}}</div>
		{{else}}
			<div><strong>No Events</strong></div>
		{{end}}
	</body>
</html>
`
)

var indexTemplate = template.Must(template.New("index").Parse(indexTemplateStr))

// Render will turn meetup event data into something to write out.
func Render(events []data.Event) []byte {
	buf := &bytes.Buffer{}
	indexTemplate.Execute(buf, events)
	return buf.Bytes()
}
