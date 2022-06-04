package plan

import (
	"time"
)

// Event is an Event with a start time and duration
type Event struct {
	Description string
	duration    time.Duration
	start       time.Time
}

func NewEvent(description string, d time.Duration) *Event {
	return &Event{
		Description: description,
		duration:    d,
	}
}

// TimeRange returns a string with start and finish time of the event, like "12:00-12:05"
func (e *Event) TimeRange() string {
	return e.timeRangeWithSep("-")
}

func (e *Event) timeRangeWithSep(sep string) string {
	beginRange := e.start.Format(Layout)
	endRange := e.start.Add(e.duration).Format(Layout)

	return beginRange + sep + endRange
}
