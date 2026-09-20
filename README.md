# cronjob
* Allow add, remove, update schedule of tasks at runtime.
* Support run with extended cron or datemath format schedule spec.
* Work with single schedule arrangement loop on a queue.

# Code Examples
```go
package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/codeindex2937/cronjob"
)

func main() {
	ctx := context.Background()
	var wg sync.WaitGroup

	cp := cronjob.NewCronParser(time.Local)
	scheduler := cronjob.NewManager[string](time.Local)

	go scheduler.Run(ctx)

	cronConfig, _ := cp.Parse("@every 1s")
	taskId := "taskID"
	var counter int32
	wg.Add(1)
	scheduler.AddTask(cronConfig, taskId, func() {
		atomic.AddInt32(&counter, 1)
		if counter > 2 {
			scheduler.RemoveTasks(taskId)
			wg.Done()
			return
		}
		fmt.Printf("count %d\n", counter)
	})

	fmt.Println("start")
	wg.Wait()
	task, _ := scheduler.FindTask(taskId)
	fmt.Printf("%v\n", task)
}
```
