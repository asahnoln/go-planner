// Package planner is designed for creating dynamic time tables.
package planner

import "time"

// Layout is a basic layout for time formatting
const Layout = "15:04"

// project presents current timetable holding events with their durations.
type project struct {
	events []*event
	start  time.Time
}

// NewProject creates a new project with default start time 12:00
func NewProject() *project {
	return &project{
		start: time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
	}
}

// AddEvent adds an event to the project with given duration.
// Use time.Duration approach to add durations, like 5 * time.Minute
func (p *project) AddEvent(d time.Duration) *event {
	e := &event{d, p.finishTime()}
	p.events = append(p.events, e)

	return e
}

// StartTime changes project's start time and shifts events start times accordingly
func (p *project) StartTime(t time.Time) {
	p.start = t
	next := t

	for _, e := range p.events {
		e.start = next
		next = next.Add(e.duration)
	}
}

// finishTime is a convinience function to get finishing time of the project
func (p *project) finishTime() time.Time {
	if l := len(p.events); l > 0 {
		lastEvent := p.events[l-1]
		return lastEvent.start.Add(lastEvent.duration)
	}

	return p.start
}
