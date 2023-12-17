package cronjob

import (
	"errors"
	"time"
)

var ErrNoItem = errors.New("key not found")

type CronQueue[T any, K comparable] struct {
	head *CronItem[T, K]
}

func newConList[T any, K comparable]() *CronQueue[T, K] {
	return &CronQueue[T, K]{head: nil}
}

func (s *CronQueue[T, K]) NextSchedule() time.Time {
	if s.head == nil {
		return time.Time{}
	}
	return s.head.sortValue
}

func (s *CronQueue[T, K]) IsEmpty() bool {
	return s.head == nil
}

func (s *CronQueue[T, K]) Insert(sortKey time.Time, key K, value T) {
	if s.head == nil {
		s.head = newCronList(sortKey, key, value)
		return
	}

	var currentNode *CronItem[T, K]
	currentNode = s.head
	var previousNode *CronItem[T, K]
	var found bool
	newNode := newCronList(sortKey, key, value)

	for {
		if currentNode.Compare(newNode) >= 0 {
			if previousNode != nil {
				newNode.next = previousNode.next
				previousNode.next = newNode
			} else {
				newNode.next = s.head
				s.head = newNode
			}
			found = true
			break
		}

		if currentNode.next == nil {
			break
		}

		previousNode = currentNode
		currentNode = currentNode.next
	}

	if !found {
		currentNode.next = newNode
	}
}

func (s *CronQueue[T, K]) Find(key K) (value T, nextSched time.Time, err error) {
	var t T
	currentNode := s.head
	for {
		if currentNode.key == key {
			return currentNode.value, currentNode.sortValue, nil
		}

		if currentNode.next == nil {
			break
		}
		currentNode = currentNode.next
	}
	return t, time.Time{}, ErrNoItem
}

func (s *CronQueue[T, K]) Remove(key K) error {
	if s.head == nil {
		return ErrNoItem
	}

	currentNode := s.head
	var previousNode *CronItem[T, K]
	for {
		if currentNode.key == key {
			if previousNode != nil {
				previousNode.next = currentNode.next
			} else {
				s.head = currentNode.next
			}
			return nil
		}

		if currentNode.next == nil {
			break
		}
		previousNode = currentNode
		currentNode = currentNode.next
	}
	return ErrNoItem
}

func (s *CronQueue[T, K]) Pop() (T, error) {
	if s.head == nil {
		var t T
		return t, ErrNoItem
	}

	popped := s.head
	s.head = popped.next
	return popped.value, nil
}

func (s *CronQueue[T, K]) DisplayAll() {
	log.Infof("")
	log.Infof("head->")
	currentNode := s.head
	for {
		log.Infof("[key:%v][val:%v]->", currentNode.sortValue, currentNode.sortValue)
		if currentNode.next == nil {
			break
		}
		currentNode = currentNode.next
	}
	log.Infof("nil\n")
}
