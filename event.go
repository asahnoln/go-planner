package planner

import (
	"time"
)

// Event is an Event with a start time and duration
type Event struct {
	description string
	duration    time.Duration
	start       time.Time
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

// FullDescription returns event's descrition with its time range
func (e *Event) FullDescription() string {
	return e.description + ": " + e.TimeRange()
}
