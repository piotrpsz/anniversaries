package events

import (
	"fmt"
)

type UserData = map[string]any

type Event struct {
	eid  EventID
	data UserData
}

func NewEventWithID(id EventID) *Event {
	return &Event{
		eid:  id,
		data: make(UserData),
	}
}

func (e *Event) Add(key string, value any) *Event {
	e.data[key] = value
	return e
}

func NewEvent(id EventID, data map[string]any) *Event {
	return &Event{
		eid:  id,
		data: data,
	}
}

func (e *Event) String() string {
	return fmt.Sprintf("%v => %v", e.eid, e.data)
}

func (e *Event) Id() EventID {
	return e.eid
}

func (e *Event) Data() map[string]any {
	return e.data
}
