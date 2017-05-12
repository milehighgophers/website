package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
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
	addr := "localhost:8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
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

	eventCh := make(chan []eventData)
	done := make(chan struct{})
	var all []eventData

	go func() {
		for eds := range eventCh {
			all = append(all, eds...)
		}
		close(done)
	}()

	var wg sync.WaitGroup
	for _, meetup := range meetupNames {
		wg.Add(1)
		go events(meetup, eventCh, &wg)
	}
	wg.Wait()
	close(eventCh)

	<-done
	return all
}

func events(name string, out chan []eventData, wg *sync.WaitGroup) {
	defer wg.Done()

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
	out <- data
}
