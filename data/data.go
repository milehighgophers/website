package data

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
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
type Store struct {
	pollingInterval time.Duration

	mu         sync.Mutex
	eventCache []Event
}

// NewStore creates a new store initialized with a polling interval.
func NewStore(i time.Duration) *Store {
	return &Store{
		pollingInterval: i,
	}
}

// Poll runs forever, polling the meetup API for event data and updating the
// internal cache.
func (s *Store) Poll() {
	for {
		events := s.poll()
		s.updateCache(events)
		time.Sleep(s.pollingInterval)
	}
}

func (s *Store) updateCache(events []Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.eventCache = events
}

func (s *Store) poll() []Event {
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

	sort.Slice(all, func(i, j int) bool {
		return all[i].Time < all[j].Time
	})
	return all
}

// AllEvents returns the current meetup events in CO.
func (s *Store) AllEvents() []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.eventCache
}

// Event contains information about a meetup event.
type Event struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Time int64  `json:"time"`
}

// HumanTime returns the time formated for the UI.
func (e Event) HumanTime() string {
	return time.Unix(e.Time/1000, 0).Format(time.UnixDate)
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
