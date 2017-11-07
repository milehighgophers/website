package data

import (
	"errors"
	"reflect"
	"testing"
)

type TestClient struct {
	*MeetupClient

	data []byte
	err  error
}

func NewClient(data string, err error) TestClient {
	return TestClient{
		data: []byte(data),
		err:  err,
	}
}

func (c TestClient) Get(key string) ([]byte, error) {
	if c.err != nil {
		return []byte{}, c.err
	}

	return c.data, c.err
}

func TestScheduleHasNoEvents(t *testing.T) {
	s := Schedule{}
	if len(s.Events) != 0 {
		t.Fail()
	}
}

func TestFetchEventsStoresEvents(t *testing.T) {
	s := Schedule{key: "Boulder-Gophers"}

	c := NewClient(`[{"id":"id","name":"Event","time":400}]`, nil)

	s.FetchEvents(c)

	expected := []Event{Event{ID: "id", Name: "Event", Time: 400}}

	if !reflect.DeepEqual(s.Events, expected) {
		t.Fail()
	}
}

func TestFetchEventsReturnsErrorWhenRequestFails(t *testing.T) {
	s := Schedule{key: "Boulder-Gophers"}
	c := NewClient(``, errors.New("Failed to connect"))

	err := s.FetchEvents(c)

	if err == nil {
		t.Fail()
	}
}

func TestFetchEventsReturnsErrorWhenUnmarshalFails(t *testing.T) {
	s := Schedule{key: "Boulder-Gophers"}
	c := NewClient(`<>`, nil)

	err := s.FetchEvents(c)

	if err == nil {
		t.Fail()
	}
}

func TestFetchEventsReturnsEventsInOrderOfTime(t *testing.T) {
	s := Schedule{key: "Boulder-Gophers"}
	c := NewClient(`
		[
			{"id":"two","name":"Two","time":2},
			{"id":"one","name":"One","time":1}
		]
	`, nil)

	s.FetchEvents(c)

	expected := []Event{
		Event{ID: "one", Name: "One", Time: 1},
		Event{ID: "two", Name: "Two", Time: 2},
	}

	if !reflect.DeepEqual(s.Events, expected) {
		t.Fail()
	}
}

func TestNextReturnsSubsetOfEvents(t *testing.T) {
	s := Schedule{Events: []Event{
		Event{ID: "one", Name: "One", Time: 1},
		Event{ID: "two", Name: "Two", Time: 2},
		Event{ID: "three", Name: "Three", Time: 3},
	}}

	expected := []Event{
		Event{ID: "one", Name: "One", Time: 1},
		Event{ID: "two", Name: "Two", Time: 2},
	}

	if !reflect.DeepEqual(s.Next(2), expected) {
		t.Fail()
	}
}

func TestNextReturnsAllWhenLimitGreaterThanLen(t *testing.T) {
	s := Schedule{Events: []Event{
		Event{ID: "one", Name: "One", Time: 1},
	}}

	expected := []Event{
		Event{ID: "one", Name: "One", Time: 1},
	}

	if !reflect.DeepEqual(s.Next(2), expected) {
		t.Fail()
	}
}

func TestHumanTimeReturnsTimeInProperFormat(t *testing.T) {
	e := Event{Time: 1533121600000}

	if e.HumanTime() != "Wed, 01 Aug 2018 05:06:40 MDT" {
		t.Fail()
	}
}
