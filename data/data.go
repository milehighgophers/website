package data

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

const apiTemplate = "https://api.meetup.com/%s/events?status=upcoming"

var (
	meetupNames = []string{
		"Boulder-Gophers",
		"Denver-Go-Language-User-Group",
		"Denver-Go-Programming-Language-Meetup",
	}
)

// Store contains data for the site.
type Store struct{}

// AllEvents returns the current meetup events in CO.
func (s *Store) AllEvents() []Event {
	// TODO: consider sorting by timestamp

	eventCh := make(chan []Event)
	done := make(chan struct{})
	var all []Event

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

// Event contains information about a meetup event.
type Event struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Time int64  `json:"time"`
}

func events(name string, out chan []Event, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(fmt.Sprintf(apiTemplate, name))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var data []Event
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	out <- data
}
