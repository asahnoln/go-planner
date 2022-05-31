// Package plan is designed for creating dynamic time tables.
package plan

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Layout is a basic layout for time formatting
const Layout = "15:04"

// Project presents current timetable holding events with their durations.
type Project struct {
	events []*Event
	start  time.Time
}

// NewProject creates a new project with default start time 12:00
func NewProject() *Project {
	return &Project{
		start: time.Date(1, 1, 1, 12, 0, 0, 0, time.UTC),
	}
}

func (p *Project) Add(e *Event) *Project {
	e.start = p.finishTime()
	p.events = append(p.events, e)
	return p
}

func (p *Project) Event(i int) *Event {
	return p.events[i]
}

// AddEvent adds an event to the project with given duration.
// Use time.Duration approach to add durations, like 5 * time.Minute
func (p *Project) AddEvent(description string, d time.Duration) *Event {
	e := &Event{
		description: description,
		duration:    d,
		start:       p.finishTime(),
	}
	p.events = append(p.events, e)

	return e
}

// StartTime changes project's start time and shifts events start times accordingly
func (p *Project) StartTime(t time.Time) {
	p.start = t
	next := t

	for _, e := range p.events {
		e.start = next
		next = next.Add(e.duration)
	}
}

// Table writes a timetable of the project to the given Writer
func (p *Project) Table(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	for _, e := range p.events {
		fmt.Fprintf(tw, "%s\t| %s\n", e.description, e.timeRangeWithSep(" | "))
	}
	tw.Flush()
}

// finishTime is a convinience function to get finishing time of the project
func (p *Project) finishTime() time.Time {
	if l := len(p.events); l > 0 {
		lastEvent := p.events[l-1]
		return lastEvent.start.Add(lastEvent.duration)
	}

	return p.start
}
