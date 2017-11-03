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

	mu             sync.Mutex
	meetupSchedule *MeetupSchedule
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

func (s *Store) updateCache(schedule *MeetupSchedule) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.meetupSchedule = schedule
}

func (s *Store) poll() *MeetupSchedule {
	schedule := NewMeetupSchedule()
	for _, meetup := range meetupNames {
		eds, err := events(meetup)
		if err != nil {
			log.Printf("error fetching events for %s: %s", meetup, err)
			continue
		}
		sort.Slice(eds, func(i, j int) bool {
			return eds[i].Time < eds[j].Time
		})
		schedule.Add(meetup, eds)
	}
	return schedule
}

// AllEvents returns the current meetup events in CO.
func (s *Store) AllEvents() *MeetupSchedule {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.meetupSchedule
}

// Event contains information about a meetup event.
type Event struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Time int64  `json:"time"`
}

// HumanTime returns the time formated for the UI.
func (e Event) HumanTime() string {
	return time.Unix(e.Time/1000, 0).Format(time.RFC1123)
}

func NewMeetupSchedule() *MeetupSchedule {
	return &MeetupSchedule{
		events: make(map[string][]Event),
	}
}

type MeetupSchedule struct {
	events map[string][]Event
}

func (m *MeetupSchedule) Add(name string, events []Event) {
	m.events[name] = events
}

func (m *MeetupSchedule) BoulderEvents() []Event {
	return nextThree(m.events["Boulder-Gophers"])
}

func (m *MeetupSchedule) DenverEvents() []Event {
	return nextThree(m.events["Denver-Go-Language-User-Group"])
}

func (m *MeetupSchedule) DTCEvents() []Event {
	return nextThree(m.events["Denver-Go-Programming-Language-Meetup"])
}

func nextThree(events []Event) []Event {
	if len(events) < 3 {
		return events
	}

	return events[0:3]
}

func events(name string) (data []Event, err error) {
	resp, err := http.Get(fmt.Sprintf(apiTemplate, name))

	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&data)

	if err != nil {
		return data, err
	}

	return data, err
}
