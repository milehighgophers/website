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
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/normalize/7.0.0/normalize.css">
		<link rel="stylesheet" href="/assets/styles.css">
	</head>
	<body>
		<header class="header">
			<img src="/assets/logo.png" alt="mile-high-gophers-logo" class="logo">
			<img src="/assets/hero.jpg" alt="mountains-backdrop">
		</header>
		<section class="body flex-container">
			<div class="flex-1">
				<h1>Boulder</h1>
				<ul>
				{{range .BoulderEvents}}
					<li>{{.HumanTime}} -- {{.Name}}</li>
				{{else}}
					<div><strong>No Events</strong></div>
				{{end}}
				</ul>
			</div>
			<div class="flex-1">
				<h1>Denver</h1>
				<ul>
				{{range .DenverEvents}}
					<li>{{.HumanTime}} -- {{.Name}}</li>
				{{else}}
					<div><strong>No Events</strong></div>
				{{end}}
				</ul>
			</div>
			<div class="flex-1">
				<h1>Denver Tech Center</h1>
				<ul>
				{{range .DTCEvents}}
					<li>{{.HumanTime}} -- {{.Name}}</li>
				{{else}}
					<div><strong>No Events</strong></div>
				{{end}}
				</ul>
			</div>
		</section>
	</body>
</html>
`
)

var indexTemplate = template.Must(template.New("index").Parse(indexTemplateStr))

// Render will turn meetup event data into something to write out.
func Render(s *data.MeetupSchedule) []byte {
	buf := &bytes.Buffer{}
	indexTemplate.Execute(buf, s)
	return buf.Bytes()
}
