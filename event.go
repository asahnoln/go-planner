package planner

import (
	"time"
)

// Event is an event with a start time and duration
type Event struct {
	duration time.Duration
	start    time.Time
}

// TimeRange returns a string with start and finish time of the event, like "12:00-12:05"
func (e *Event) TimeRange() string {
	beginRange := e.start.Format(Layout)
	endRange := e.start.Add(e.duration).Format(Layout)

	return beginRange + "-" + endRange
}
