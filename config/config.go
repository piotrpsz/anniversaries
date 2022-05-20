package config

import (
	"image/color"
	"sync"

	events "calendar/event"
)

var (
	instance *Settings
	once     sync.Once
)

type Settings struct {
	darkTheme bool
}

func Config() *Settings {
	once.Do(func() {
		instance = newConfig()
	})
	return instance
}

func newConfig() *Settings {
	return &Settings{
		darkTheme: true,
	}
}

func (s *Settings) DarkTheme() bool {
	return s.darkTheme
}

func (s *Settings) SetDarkTheme(state bool) {
	s.darkTheme = state
	changeTheme(state)
	e := events.NewEventWithID(events.RefreshRequest)
	events.Instance().Send(e)
}

func (s *Settings) ChangeTheme() {
	if s.darkTheme {
		s.darkTheme = false
	} else {
		s.darkTheme = true
	}
	changeTheme(s.darkTheme)
	e := events.NewEventWithID(events.RefreshRequest)
	events.Instance().Send(e)
}

func (s *Settings) TextColor() color.NRGBA {
	if s.darkTheme {
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	}
	return color.NRGBA{A: 0xff}
}
