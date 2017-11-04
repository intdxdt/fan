package fan

import (
	"sync"
)

//process
func process(worker func(interface{}) interface{},
	in <-chan interface{}, concurrency int, exit <-chan struct{}) <-chan interface{} {
	var wg sync.WaitGroup
	//set up number of of clones to wait for
	wg.Add(concurrency)
	var out = make(chan interface{}, 2*concurrency)
	var onExit = false
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
