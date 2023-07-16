package fan

import (
	"sync"
)

func process[T any](worker func(T) T, in <-chan T, concurrency int, exit <-chan struct{}) <-chan T {
	var wg sync.WaitGroup
	wg.Add(concurrency)
	var out = make(chan T, 2*concurrency)
	var onExit = false

	var fn = func(idx int) {
		defer wg.Done()
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

	//wait for all the clones to be done in a new go routine
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
