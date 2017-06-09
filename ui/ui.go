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
	{{range $key, $value := .}}
		<h1>{{$key}}</h1>
		<ul>
		{{range $value}}
			<li>{{.HumanTime}} -- {{.Name}}</li>
		{{else}}
			<div><strong>No Events</strong></div>
		{{end}}
		</ul>
	{{end}}
	</body>
</html>
`
)

var indexTemplate = template.Must(template.New("index").Parse(indexTemplateStr))

// Render will turn meetup event data into something to write out.
func Render(events map[string][]data.Event) []byte {
	buf := &bytes.Buffer{}
	indexTemplate.Execute(buf, events)
	return buf.Bytes()
}
