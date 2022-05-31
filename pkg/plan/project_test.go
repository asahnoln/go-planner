package plan_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/asahnoln/go-planner/pkg/plan"
	"github.com/stretchr/testify/assert"
)

func TestAddEventGetTimeRange(t *testing.T) {
	p := plan.NewProject()

	event := plan.NewEvent("Introduction", 5*time.Minute)
	event2 := plan.NewEvent("Warmup", 10*time.Minute)
	event3 := plan.NewEvent("Zip Zap Zop", 15*time.Minute)

	p.Add(event).
		Add(event2).
		Add(event3)

	assert.Equal(t, "12:15-12:30", p.Event(2).TimeRange())
	assert.Equal(t, "12:00-12:05", p.Event(0).TimeRange())
	assert.Equal(t, "12:05-12:15", p.Event(1).TimeRange())
}

func TestChangeProjectTimeAfterEvent(t *testing.T) {
	p := plan.NewProject()
	e := p.AddEvent("Warmup", 5*time.Minute)
	p.StartTime(time.Date(2000, 1, 1, 15, 30, 0, 0, time.UTC))

	assert.Equal(t, "15:30-15:35", e.TimeRange())
}

func TestChangeProjectTimeBeforeEvent(t *testing.T) {
	p := plan.NewProject()
	p.StartTime(time.Date(2000, 1, 1, 15, 30, 0, 0, time.UTC))
	e := p.AddEvent("Warmup", 5*time.Minute)

	assert.Equal(t, "15:30-15:35", e.TimeRange())
}

func TestProjectTable(t *testing.T) {
	p := plan.NewProject()
	p.Add(plan.NewEvent("Intro", 5*time.Minute))
	p.Add(plan.NewEvent("Warmup", 10*time.Minute))
	p.Add(plan.NewEvent("Zip Zap Zop", 15*time.Minute))

	var buf bytes.Buffer
	p.Table(&buf)

	want := `
Intro       | 12:00 | 12:05
Warmup      | 12:05 | 12:15
Zip Zap Zop | 12:15 | 12:30
`[1:]

	assert.Equal(t, want, buf.String())
}
