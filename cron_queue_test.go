package cronjob

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLinkedNodeCompare(t *testing.T) {
	as := assert.New(t)
	now := time.Now()
	n1 := newCronList(now, "1", "node1")
	n2 := newCronList(now.Add(1*time.Second), "2", "node2")

	as.Equal(-1, n1.Compare(n2))
	as.Equal(1, n2.Compare(n1))
	as.Equal(0, n1.Compare(n1))
}

func TestScheduleListSearch(t *testing.T) {
	as := assert.New(t)
	now := time.Now()
	targetSched := now.Add(4 * time.Second)
	schedList := newConList[string, string]()
	schedList.Insert(now.Add(3*time.Second), "3", "string3")
	schedList.Insert(targetSched, "4", "string4")
	schedList.Insert(now.Add(2*time.Second), "2", "string2")

	ret, nextSched, err := schedList.Find("4")
	as.NoError(err)
	as.Equal("string4", ret)
	as.Equal(nextSched, targetSched)
}

func TestScheduleListPop(t *testing.T) {
	as := assert.New(t)
	now := time.Now()
	schedList := newConList[string, string]()
	schedList.Insert(now.Add(3*time.Second), "3", "string3")
	schedList.Insert(now.Add(4*time.Second), "4", "string4")
	schedList.Insert(now.Add(2*time.Second), "2", "string2")

	ret, err := schedList.Pop()
	as.NoError(err)
	as.Equal("string2", ret)
	ret, err = schedList.Pop()
	as.NoError(err)
	as.Equal("string3", ret)
	ret, err = schedList.Pop()
	as.NoError(err)
	as.Equal("string4", ret)
	ret, err = schedList.Pop()
	as.Equal(ErrNoItem, err)
	as.Equal("", ret)
}

func TestScheduleListRemove(t *testing.T) {
	as := assert.New(t)
	now := time.Now()
	schedList := newConList[string, string]()
	schedList.Insert(now.Add(4*time.Second), "4", "string4")
	schedList.Insert(now.Add(3*time.Second), "3", "string3")
	as.Nil(schedList.Remove("3"))
	as.Nil(schedList.Remove("4"))
	as.Equal(ErrNoItem, schedList.Remove("4"))
}

func TestScheduleListIsEmpty(t *testing.T) {
	as := assert.New(t)
	now := time.Now()
	schedList := newConList[string, string]()
	as.True(schedList.IsEmpty())

	schedList.Insert(now.Add(4*time.Second), "4", "string4")
	as.False(schedList.IsEmpty())

	value, err := schedList.Pop()
	as.NoError(err)
	as.True(schedList.IsEmpty())
	as.Equal("string4", value)

	schedList.Insert(now.Add(3*time.Second), "3", "string3")
	as.False(schedList.IsEmpty())

	err = schedList.Remove("3")
	as.NoError(err)
	as.True(schedList.IsEmpty())
}
