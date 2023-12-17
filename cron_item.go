package cronjob

import "time"

type CronItem[T any, K comparable] struct {
	sortValue time.Time
	key       K
	value     T
	next      *CronItem[T, K]
}

func newCronList[T any, K comparable](sortValue time.Time, key K, value T) *CronItem[T, K] {
	return &CronItem[T, K]{sortValue, key, value, nil}
}

func (node *CronItem[T, K]) Compare(that *CronItem[T, K]) int {
	if node.sortValue.Before(that.sortValue) {
		return -1
	}
	if node.sortValue.After(that.sortValue) {
		return 1
	}
	return 0
}

func (node *CronItem[T, K]) Value() *T {
	return &node.value
}
