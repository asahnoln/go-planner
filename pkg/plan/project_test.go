package plan_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/asahnoln/go-planner/pkg/plan"
	"github.com/stretchr/testify/assert"
)

func TestAddEventGetTimeRange(t *testing.T) {
	p := plan.New()

	es := p.Add(plan.NewEvent("Introduction", 5*time.Minute),
		plan.NewEvent("Warmup", 10*time.Minute),
		plan.NewEvent("Zip Zap Zop", 15*time.Minute))

	assert.Equal(t, "12:15-12:30", es[2].TimeRange())
	assert.Equal(t, "12:00-12:05", es[0].TimeRange())
	assert.Equal(t, "12:05-12:15", es[1].TimeRange())
}

func TestChangeProjectTimeAfterEvent(t *testing.T) {
	p := plan.New()
	es := p.Add(plan.NewEvent("Warmup", 5*time.Minute))
	p.StartTime(time.Date(2000, 1, 1, 15, 30, 0, 0, time.UTC))

	assert.Equal(t, "15:30-15:35", es[0].TimeRange())
}

func TestChangeProjectTimeBeforeEvent(t *testing.T) {
	p := plan.New()
	p.StartTime(time.Date(2000, 1, 1, 15, 30, 0, 0, time.UTC))
	es := p.Add(plan.NewEvent("Warmup", 5*time.Minute))

	assert.Equal(t, "15:30-15:35", es[0].TimeRange())
}

func TestProjectTable(t *testing.T) {
	p := plan.New()
	p.Add(plan.NewEvent("Intro", 5*time.Minute),
		plan.NewEvent("Warmup", 10*time.Minute),
		plan.NewEvent("Zip Zap Zop", 15*time.Minute))

	var buf bytes.Buffer
	p.Table(&buf)

	want := `
Intro       | 12:00 | 12:05
Warmup      | 12:05 | 12:15
Zip Zap Zop | 12:15 | 12:30
`[1:]

	assert.Equal(t, want, buf.String())
}
