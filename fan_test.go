package fan

import (
	"fmt"
	"github.com/franela/goblin"
	"golang.org/x/exp/constraints"
	"runtime"
	"testing"
	"time"
)

var data = []int{40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40}

func slowFib[T constraints.Signed](n T) T {
	if n == 0 || n == 1 {
		return n
	}
	return slowFib(n-1) + slowFib(n-2)
}

func oneByOne() []int {
	var results = make([]int, len(data))
	for i, d := range data {
		results[i] = slowFib(d)
	}
	return results
}

func streamer() []int {
	var results = make([]int, 0)
	var stream = make(chan int)
	var exit = make(chan struct{})
	defer close(exit)

	go func() {
		for _, d := range data {
			stream <- d
		}
		close(stream)
	}()

	var worker = func(v int) int {
		return slowFib(v)
	}

	var out = Stream(stream, worker, runtime.NumCPU(), exit)
	for o := range out {
		results = append(results, o)
	}
	return results
}

func payload() []int {
	var exit = make(chan struct{})
	defer close(exit)

	var worker = func(v int) int {
		return slowFib(v)
	}

	return Payload(data, worker, runtime.NumCPU(), exit)
}

func timeIt[T any](fn func() []T, desc string) []T {
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

			var expects = timeIt(oneByOne, "Serial")

			var strm = timeIt(streamer, "Streamer")
			g.Assert(strm).Equal(expects)

			var pyld = timeIt(payload, "Payload")
			g.Assert(pyld).Equal(expects)

		})
	})
}
