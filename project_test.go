package planner_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/asahnoln/go-planner"
)

func TestAddEventGetTimeRange(t *testing.T) {
	p := planner.NewProject()

	event := p.AddEvent("Introduction", 5*time.Minute)
	event2 := p.AddEvent("Warmup", 10*time.Minute)
	event3 := p.AddEvent("Zip Zap Zop", 15*time.Minute)

	message := "want time range %q, got %q"
	assertSameString(t, "12:15-12:30", event3.TimeRange(), message)
	assertSameString(t, "12:00-12:05", event.TimeRange(), message)
	assertSameString(t, "12:05-12:15", event2.TimeRange(), message)
}

func TestChangeProjectTimeAfterEvent(t *testing.T) {
	p := planner.NewProject()
	e := p.AddEvent("Warmup", 5*time.Minute)
	p.StartTime(time.Date(2000, 1, 1, 15, 30, 0, 0, time.UTC))

	assertSameString(t, "15:30-15:35", e.TimeRange(), "want time range with changed base %q, got %q")
}

func TestChangeProjectTimeBeforeEvent(t *testing.T) {
	p := planner.NewProject()
	p.StartTime(time.Date(2000, 1, 1, 15, 30, 0, 0, time.UTC))
	e := p.AddEvent("Warmup", 5*time.Minute)

	assertSameString(t, "15:30-15:35", e.TimeRange(), "want time range with changed base %q, got %q")
}

func TestProjectTable(t *testing.T) {
	p := planner.NewProject()
	p.AddEvent("Intro", 5*time.Minute)
	p.AddEvent("Warmup", 10*time.Minute)
	p.AddEvent("Zip Zap Zop", 15*time.Minute)

	var buf bytes.Buffer
	p.Table(&buf)

	want := `
Intro       | 12:00 | 12:05
Warmup      | 12:05 | 12:15
Zip Zap Zop | 12:15 | 12:30
`[1:]

	assertSameString(t, want, buf.String(), "want table\n%v\n\ngot\n%v\n")
}

func assertSameString(t testing.TB, want, got, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}
