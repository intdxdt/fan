package fan

import (
	"github.com/franela/goblin"
	"time"
	"fmt"
	"testing"
	"runtime"
)

const ConCur = 8

var data = []int{40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40}

func slowFib(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return slowFib(n-1) + slowFib(n-2)
}

func onbyone() []interface{} {
	var results = make([]interface{}, len(data))
	for i, d := range data {
		results[i] = slowFib(d)
	}
	return results
}

func streamer() []interface{} {
	var results = make([]interface{}, 0)
	var stream = make(chan interface{})
	var exit = make(chan struct{})
	defer close(exit)

	go func() {
		for _, d := range data {
			stream <- d
		}
		close(stream)
	}()

	var worker = func(v interface{}) interface{} {
		return slowFib(v.(int))
	}

	var out = Stream(stream, worker, ConCur, exit)
	for o := range out {
		results = append(results, o.(int))
	}
	return results
}

func payload() []interface{} {
	var dat = make([]interface{}, 0)
	var exit = make(chan struct{})
	defer close(exit)

	for _, d := range data {
		dat = append(dat, d)
	}

	worker := func(v interface{}) interface{} {
		return slowFib(v.(int))
	}

	res := Payload(dat, worker, ConCur, exit)
	return res
}

func pool() []interface{} {
	var results = make([]interface{}, 0)
	var stream = make(chan interface{})
	var exit = make(chan struct{})

	go func() {
		for _, d := range data {
			stream <- d
		}
	}()

	var worker = func(v interface{}) interface{} {
		return slowFib(v.(int))
	}

	var done = make(chan struct{})
	var pool = NewPool(stream, worker, ConCur, exit)

	var out = pool.Start()
	go func() {
		for {
			select {
			case o := <-out:
				results = append(results, o.(int))
			default:
				//halting condition all data served and processed
				if len(results) == len(data) && pool.IsIdle() {
					close(exit)
					close(done)
					return
				}
			}
		}
	}()
	<-done
	return results
}

func timeIt(fn func() []interface{}, desc string) interface{} {
	var tick = time.Now()
	var res  = fn()
	var tock = time.Now()
	var d    = tock.Sub(tick)
	fmt.Println(fmt.Sprintf("%v:\tduration:\t%.2f", desc, d.Seconds()))
	return res
}

func TestFan(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	g := goblin.Goblin(t)
	g.Describe("fan in - fan out", func() {
		g.It("should test serial and parallel fan", func() {
			g.Timeout(1 * time.Minute)

			expects := timeIt(onbyone, "Serial")

			strm := timeIt(streamer, "Streamer")
			g.Assert(strm).Equal(expects)

			pyld := timeIt(payload, "Payload")
			g.Assert(pyld).Equal(expects)

			pool := timeIt(payload, "Pool")
			g.Assert(pool).Equal(expects)
		})
	})
}