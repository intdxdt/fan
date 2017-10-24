package main

import (
	"fmt"
	"fan"
	"time"
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

	var out = fan.IOStream(stream, worker, ConCur, exit)
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

	res := fan.Payload(dat, worker, ConCur, exit)
	return res
}

func timeIt(fn func() []interface{}, desc string) interface{} {
	tick := time.Now()
	res := fn()
	tock := time.Now()
	d := tock.Sub(tick)
	fmt.Println(fmt.Sprintf("%v: , duration: %v", desc, d.Seconds()))
	return res
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	oneby := timeIt(onbyone, "OneByOne")
	fmt.Println(oneby)

	strm := timeIt(streamer, "Streamer")
	fmt.Println(strm)

	pyld := timeIt(payload, "Payload")
	fmt.Println(pyld)
}
