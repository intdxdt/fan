package fan

import (
	"fmt"
	"github.com/franela/goblin"
	"runtime"
	"testing"
	"time"
)

var data = []int{40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40}

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

	var out = Stream(stream, worker, runtime.NumCPU(), exit)
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

	return Payload(dat, worker, runtime.NumCPU(), exit)
}

func timeIt(fn func() []interface{}, desc string) interface{} {
	var tick = time.Now()
	var res = fn()
	var tock = time.Now()
	var d = tock.Sub(tick)
	fmt.Println(fmt.Sprintf("%v:\tduration:\t%.2f", desc, d.Seconds()))
	return res
}

func TestFan(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	g := goblin.Goblin(t)
	g.Describe("fan in - fan out", func() {
		g.It("should test serial and parallel fan", func() {
			g.Timeout(3 * time.Minute)

			expects := timeIt(onbyone, "Serial")

			strm := timeIt(streamer, "Streamer")
			g.Assert(strm).Equal(expects)

			pyld := timeIt(payload, "Payload")
			g.Assert(pyld).Equal(expects)

		})
	})
}
