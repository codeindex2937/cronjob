package cronjob

import (
	"errors"
	"fmt"
	"time"

	"github.com/robfig/cron"
	"github.com/timberio/go-datemath"
)

var ErrUnknownFormat = errors.New("unknown format")

type CronParser struct {
	loc *time.Location
}

type CronConfig struct {
	EndTime    time.Time
	Expression string
	cron       cron.Schedule
}

func (c CronConfig) Next(offset time.Time) time.Time {
	next := c.cron.Next(offset)
	if !c.EndTime.IsZero() && next.After(c.EndTime) {
		return time.Time{}
	}

	return next
}

func NewCronParser(loc *time.Location) *CronParser {
	return &CronParser{
		loc: loc,
	}
}

func (ts CronParser) LocalTime() time.Time {
	return time.Now().In(ts.loc)
}

func (ts CronParser) LocalTimeString(t time.Time) string {
	return t.In(ts.loc).String()
}

func (ts CronParser) Parse(s string) (CronConfig, error) {
	if t, err := datemath.ParseAndEvaluate(s, datemath.WithLocation(ts.loc)); err == nil {
		local := t.Local()
		config := fmt.Sprintf("%v %v %v %v %d *", local.Second(), local.Minute(), local.Hour(), local.Day(), local.Month())
		sched, err := cron.Parse(config)
		return CronConfig{
			EndTime:    t,
			Expression: config,
			cron:       sched,
		}, err
	}

	if cronSched, err := cron.Parse(s); err == nil {
		return CronConfig{
			cron:       cronSched,
			Expression: s,
		}, nil
	}

	return CronConfig{}, ErrUnknownFormat
}
