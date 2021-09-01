package planner

import "time"

var events []*Event

var startTime = time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
var lastTime = startTime

const Layout = "15:04"

type Event struct {
	duration time.Duration
	start    time.Time
}

func AddEvent(d time.Duration) *Event {
	e := &Event{d, lastTime}
	lastTime = lastTime.Add(d)

	return e
}

func (e *Event) TimeRange() string {
	beginRange := e.start.Format(Layout)
	endRange := e.start.Add(e.duration).Format(Layout)

	return beginRange + "-" + endRange
}
