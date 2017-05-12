package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	apiTemplate      = "https://api.meetup.com/%s/events?status=upcoming"
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

var (
	meetupNames = []string{
		"Boulder-Gophers",
		"Denver-Go-Language-User-Group",
		"Denver-Go-Programming-Language-Meetup",
	}
	indexTemplate = template.Must(template.New("index").Parse(indexTemplateStr))
)

type eventData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Time int64  `json:"time"`
}

func main() {
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func indexHandler(rw http.ResponseWriter, r *http.Request) {
	html := render(allEvents())
	_, err := rw.Write(html)
	if err != nil {
		log.Print("error occured with /:", err)
	}
}

func render(events []eventData) []byte {
	buf := &bytes.Buffer{}
	indexTemplate.Execute(buf, events)
	return buf.Bytes()
}

func allEvents() []eventData {
	// TODO: consider sorting by timestamp
	// TODO: refactor as range over slice
	boulder := events(meetupNames[0])
	denver := events(meetupNames[1])
	dtc := events(meetupNames[2])

	var allEvents []eventData
	allEvents = append(allEvents, boulder...)
	allEvents = append(allEvents, denver...)
	allEvents = append(allEvents, dtc...)
	return allEvents
}

func events(name string) []eventData {
	resp, err := http.Get(fmt.Sprintf(apiTemplate, name))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var data []eventData
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
