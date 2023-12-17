package cronjob

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var cp = NewCronParser(time.Local)

func setupTest(wg *sync.WaitGroup, s *Manager[uint64], done chan struct{}) {
	wg.Add(1)
	go func() {
		var running = true
		for running {
			select {
			case <-done:
				running = false
			case <-s.reschedule:
			}
		}
		wg.Done()
	}()
}

func addTask(s *Manager[uint64], ID uint64, minute int) {
	sched, _ := cp.Parse(fmt.Sprintf("%v * * * *", minute))
	s.AddTask(
		sched,
		ID,
		func() {},
	)
}

func TestScheduleServiceAddAndRemoveTask(t *testing.T) {
	as := assert.New(t)
	done := make(chan struct{})
	var wg sync.WaitGroup

	s := NewManager[uint64](time.Local)
	as.True(s.q.NextSchedule().IsZero())

	setupTest(&wg, s, done)

	minute := time.Now().Minute()
	addTask(s, 2, minute+2)
	as.Less(1*time.Minute, time.Until(s.q.NextSchedule()))
	as.Greater(2*time.Minute, time.Until(s.q.NextSchedule()))

	addTask(s, 1, minute+1)
	as.Less(0*time.Minute, time.Until(s.q.NextSchedule()))
	as.Greater(1*time.Minute, time.Until(s.q.NextSchedule()))

	s.RemoveTasks(1)
	as.Less(1*time.Minute, time.Until(s.q.NextSchedule()))
	as.Greater(2*time.Minute, time.Until(s.q.NextSchedule()))

	close(done)
	wg.Wait()
}

func TestScheduleServiceSearchTask(t *testing.T) {
	as := assert.New(t)
	done := make(chan struct{})
	var wg sync.WaitGroup

	s := NewManager[uint64](time.Local)

	setupTest(&wg, s, done)

	now := time.Now()
	minute := now.Minute()
	addTask(s, 1, minute+1)
	addTask(s, 2, minute+2)

	task, nextSched := s.FindTask(1)
	if as.NotNil(task) {
		as.Equal(task.ID, uint64(1))
		as.Equal(minute+1, nextSched.Minute())
	}

	task, nextSched = s.FindTask(2)
	if as.NotNil(task) {
		as.Equal(task.ID, uint64(2))
		as.Equal(minute+2, nextSched.Minute())
	}

	task, nextSched = s.FindTask(3)
	as.Nil(task)
}

func TestScheduleOnceTask(t *testing.T) {
	as := assert.New(t)
	done := make(chan struct{})
	exec := make(chan struct{})
	var wg sync.WaitGroup
	var isTaskExecuted bool

	s := NewManager[uint64](time.Local)

	wg.Add(1)
	go func() {
		var running = true
		for running {
			select {
			case <-done:
				running = false
			case <-s.reschedule:
			}
		}
		wg.Done()
	}()

	now := time.Now().Add(-1 * time.Minute).Format("2006-01-02T15:04:05")

	sched, err := cp.Parse(now)
	if !as.NoError(err) {
		return
	}

	s.AddTask(
		sched,
		1,
		func() {
			isTaskExecuted = true
			close(exec)
		},
	)

	s.runOverdueTasks()

	<-exec
	as.True(isTaskExecuted)

	close(done)
	wg.Wait()
}
