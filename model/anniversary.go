package model

import (
	"time"
)

type Anniversary struct {
	Id              uint64
	Date            time.Time
	OneDayReminder  bool
	TwoDaysReminder bool
	OneWeekReminder bool
}

func NewAnniversary(date time.Time) *Anniversary {
	return &Anniversary{
		Date:            date,
		OneDayReminder:  true,
		TwoDaysReminder: true,
		OneWeekReminder: true,
	}
}
