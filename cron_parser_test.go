package cronjob

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseOnce(t *testing.T) {
	as := assert.New(t)
	sched, err := cp.Parse("2021-12-07T01:02:03")
	as.NoError(err)
	as.Equal("2 1 7 12 *", sched.Expression)
	as.False(sched.EndTime.IsZero())
	as.Equal(2021, sched.EndTime.Year())
	as.Equal(time.Month(12), sched.EndTime.Month())
	as.Equal(7, sched.EndTime.Day())
	as.Equal(1, sched.EndTime.Hour())
	as.Equal(2, sched.EndTime.Minute())
	as.Equal(3, sched.EndTime.Second())
	as.True(sched.Next(time.Now()).IsZero())
}

func TestParseCron(t *testing.T) {
	as := assert.New(t)
	now := time.Now()
	sched, err := cp.Parse("* * * * *")
	as.NoError(err)
	as.True(sched.EndTime.IsZero())
	as.Equal("* * * * *", sched.Expression)
	as.Greater(time.Duration(0), now.Sub(sched.cron.Next(now)))
	as.Less(time.Duration(0), now.Add(time.Minute).Sub(sched.cron.Next(now)))
}
