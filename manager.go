package cronjob

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Task[K comparable] struct {
	Config        CronConfig
	ID            K
	Run           func()
	scheduledTime time.Time
}

type Manager[K comparable] struct {
	m          sync.Mutex
	q          CronQueue[*Task[K], K]
	p          *CronParser
	reschedule chan bool
}

func NewManager[K comparable](loc *time.Location) *Manager[K] {
	return &Manager[K]{
		p:          NewCronParser(loc),
		reschedule: make(chan bool),
	}
}

func (s *Manager[K]) AddTask(config CronConfig, id K, job func()) {
	s.m.Lock()
	defer s.m.Unlock()

	task := &Task[K]{
		Config:        config,
		ID:            id,
		Run:           job,
		scheduledTime: config.Next(s.p.LocalTime()),
	}
	s.q.Insert(task.scheduledTime, task.ID, task)
	s.reschedule <- true
}

func (s *Manager[K]) RemoveTasks(ids ...K) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, ID := range ids {
		_ = s.q.Remove(ID)
	}
	s.reschedule <- true
}

func (s *Manager[K]) FindTask(id K) (*Task[K], time.Time) {
	s.m.Lock()
	defer s.m.Unlock()

	task, nextSched, err := s.q.Find(id)
	if err != nil {
		if err == ErrNoItem {
			return nil, time.Time{}
		} else {
			log.Errorf("search task: %v", err)
		}
	}

	return task, nextSched
}

func (s *Manager[K]) Run(ctx context.Context) {
	running := true

	for running {
		nextSchedule := s.q.NextSchedule()
		for running && nextSchedule.IsZero() {
			// no task, wait any update
			select {
			case <-s.reschedule:
				nextSchedule = s.q.NextSchedule()
			case <-ctx.Done():
				running = false
			}
		}

		if !running {
			break
		}

		select {
		case <-s.reschedule:
		case <-ctx.Done():
			running = false
		case <-time.After(time.Until(nextSchedule)):
			s.runOverdueTasks()
		}
	}
}

func (s *Manager[K]) runOverdueTasks() {
	nextTimeSlot := s.p.LocalTime()
	overdueTasks, _ := s.removeOverdueTask()

	for _, task := range overdueTasks {
		nextSchedule := task.Config.Next(nextTimeSlot)
		if task.scheduledTime.After(nextSchedule) {
			// once task
			continue
		} else {
			task.scheduledTime = nextSchedule
			s.q.Insert(nextSchedule, task.ID, task)
		}
	}

	for i := range overdueTasks {
		task := overdueTasks[i]
		go func() {
			task.Run()
		}()
	}
}

func (s *Manager[K]) removeOverdueTask() (overdueTasks []*Task[K], err error) {
	overdueTasks = []*Task[K]{}

	for !s.q.IsEmpty() && time.Until(s.q.NextSchedule()) <= 0 {
		value, err := s.q.Pop()
		if err != nil {
			log.Error("unexpected nil task\n")
			break
		}

		overdueTasks = append(overdueTasks, value)
	}
	return
}
