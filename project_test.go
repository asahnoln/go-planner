package planner_test

import (
	"testing"
	"time"

	"github.com/asahnoln/go-planner"
)

func TestAddEventGetTimeRange(t *testing.T) {
	project := planner.NewProject()

	event := project.AddEvent("Introduction", 5*time.Minute)
	event2 := project.AddEvent("Warmup", 10*time.Minute)
	event3 := project.AddEvent("Zip Zap Zop", 15*time.Minute)

	message := "want time range %q, got %q"
	assertSameString(t, "12:15-12:30", event3.TimeRange(), message)
	assertSameString(t, "12:00-12:05", event.TimeRange(), message)
	assertSameString(t, "12:05-12:15", event2.TimeRange(), message)
}

func TestChangeProjectTimeAfterEvent(t *testing.T) {
	project := planner.NewProject()
	event := project.AddEvent("Warmup", 5*time.Minute)
	project.StartTime(time.Date(2000, 1, 1, 15, 30, 0, 0, time.UTC))

	assertSameString(t, "15:30-15:35", event.TimeRange(), "want time range with changed base %q, got %q")
}

func TestChangeProjectTimeBeforeEvent(t *testing.T) {
	project := planner.NewProject()
	project.StartTime(time.Date(2000, 1, 1, 15, 30, 0, 0, time.UTC))
	event := project.AddEvent("Warmup", 5*time.Minute)

	assertSameString(t, "15:30-15:35", event.TimeRange(), "want time range with changed base %q, got %q")
}

func assertSameString(t testing.TB, want, got, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}
