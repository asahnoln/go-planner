package plan

import (
	"time"
)

// Event is an Event with duration
type Event struct {
	Description string
	duration    time.Duration
}

type ProjectEvent struct {
	*Event
	start   time.Time
	project *Project
}

func NewEvent(description string, d time.Duration) *Event {
	return &Event{
		Description: description,
		duration:    d,
	}
}

func (e *Event) Duration() time.Duration {
	return e.duration
}

// TimeRange returns a string with start and finish time of the event, like "12:00-12:05"
func (e *ProjectEvent) TimeRange() string {
	return e.timeRangeWithSep("-")
}

func (e *ProjectEvent) timeRangeWithSep(sep string) string {
	beginRange := e.start.Format(Layout)
	endRange := e.start.Add(e.duration).Format(Layout)

	return beginRange + sep + endRange
}

func (e *ProjectEvent) SetDuration(d time.Duration) {
	e.duration = d
	e.project.StartTime(e.project.start)
}
