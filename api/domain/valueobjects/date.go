package valueobjects

import (
	"time"
)

type Date struct {
	Time time.Time
}

func NewDate(t time.Time) Date {
	return Date{Time: t}
}
