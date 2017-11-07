package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Client interface {
	Get(string) ([]byte, error)
}

type MeetupClient struct{}

func (c MeetupClient) Get(url string) (data []byte, err error) {
	resp, err := http.Get(url)

	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)

	return data, err
}

type Event struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Time int64  `json:"time"`
}

// HumanTime returns the time formated for the UI.
func (e Event) HumanTime() string {
	return time.Unix(e.Time/1000, 0).Format(time.RFC1123)
}

type Schedule struct {
	key string

	Label  string
	Events []Event
}

type Schedules []*Schedule

func (s *Schedule) FetchEvents(client Client) (err error) {
	url := fmt.Sprintf("https://api.meetup.com/%s/events?status=upcoming", s.key)

	data, err := client.Get(url)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &s.Events)

	if err != nil {
		return err
	}

	sort.SliceStable(s.Events, func(i, j int) bool {
		return s.Events[i].Time < s.Events[j].Time
	})

	return err
}

func (s *Schedule) Next(count int) []Event {
	if len(s.Events) > count {
		return s.Events[0:count]
	}

	return s.Events
}

// Store contains data for the site.
type Store struct {
	pollingInterval time.Duration
	mu              sync.Mutex

	Schedules Schedules
}

// NewStore creates a new store initialized with a polling interval.
func NewStore(i time.Duration) *Store {
	return &Store{
		pollingInterval: i,
		Schedules: Schedules{
			&Schedule{key: "Boulder-Gophers", Label: "Boulder"},
			&Schedule{key: "Denver-Go-Language-User-Group", Label: "Denver"},
			&Schedule{key: "Denver-Go-Programming-Language-Meetup", Label: "Denver Tech Center"},
		},
	}
}

// Poll runs forever, polling the meetup API for event data and updating the
// internal cache.
func (s *Store) Poll() {
	for {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.refresh()

		time.Sleep(s.pollingInterval)
	}
}

func (s *Store) refresh() {
	client := MeetupClient{}

	for _, s := range s.Schedules {
		err := s.FetchEvents(client)

		if err != nil {
			log.Printf("error fetching events for %s: %s", s.key, err)
			continue
		}
	}
}
