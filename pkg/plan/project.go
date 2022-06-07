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
	events []*ProjectEvent
	start  time.Time
}

// New creates a new project with default start time 12:00
func New() *Project {
	return &Project{
		start: time.Date(1, 1, 1, 12, 0, 0, 0, time.UTC),
	}
}

func (p *Project) Add(es ...*Event) []*ProjectEvent {
	for _, e := range es {
		pe := &ProjectEvent{
			Event:   e,
			start:   p.finishTime(),
			project: p,
		}
		p.events = append(p.events, pe)
	}

	return p.events
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
		fmt.Fprintf(tw, "%s\t| %s\n", e.Description, e.timeRangeWithSep(" | "))
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
