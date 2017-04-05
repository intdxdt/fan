package fan

import (
	"sync"
	"github.com/matryer/runner"
)

//process
func process(worker func(interface{}) interface{},
	in <-chan interface{}, concurrency int, exit <-chan struct{}) <-chan interface{} {
	var wg sync.WaitGroup
	//set up number of of clones to wait for
	wg.Add(concurrency)
	out := make(chan interface{})
	onExit := false
	//assume only one worker reading from input chan
	fn := func(idx int) {
		defer wg.Done()
		//perform fn here...
		for o := range in {
			select {
			case <-exit:
				onExit = true
				return
			default:
				if onExit {
					return
				}
				out <- worker(o)
			}
		}
	}

	//now expand one worker into clones of workers
	go func() {
		for i := 0; i < concurrency; i++ {
			go fn(i)
		}
	}()

	//wait for all the clones to be done
	//in a new go routine
	go func() {
		wg.Wait()
		close(out)
	}()

	//return out chan to whoever want to read from it
	return out
}

func processTask(in <-chan interface{}, worker func(interface{}) interface{}) (*runner.Task, <-chan interface{}) {
	out := make(chan interface{})
	task := runner.Go(func(onStop runner.S) error {
		defer close(out)
		for val := range in {
			if onStop() {
				break
			}
			out <- worker(val)
		}
		return nil
	})
	return task, out
}

func processTaskPool(worker func(interface{}) interface{},
	in <-chan interface{}, concurrency int, exit <-chan struct{}) <-chan interface{} {
	var wg sync.WaitGroup
	//set up number of of clones to wait for
	wg.Add(concurrency)
	out := make(chan interface{})
	onExit := false

	//task pool
	pool := make([]*runner.Task, 0)

	//assume only one worker reading from input chan
	tasker := func(id int) *runner.Task {
		var task *runner.Task

		stop_task := func() {
			if task != nil {
				task.Stop()
			}
		}

		fn := func(onStop runner.S) error {
			defer wg.Done()
			//perform fn here...
		loop:
			for val := range in {
				select {
				case <-exit:
					onExit = true
					stop_task()
					break loop
				default:
					if onExit || onStop() {
						stop_task()
						break loop
					}
					out <- worker(val)
				}
			}
			return nil
		}
		task = runner.Go(fn)
		return task
	}

	//now expand into clones of workers
	go func() {
		for i := 0; i < concurrency; i++ {
			pool = append(pool, tasker(i))
		}
	}()

	//wait for all the clones to be done
	//in a new go routine
	go func() {
		wg.Wait()
		for _, tsk := range pool {
			//signal all running task to stop
			if tsk.Running() {
				tsk.Stop()
			}
		}
		close(out)
	}()

	//return out chan to whoever want to read from it
	return out
}
