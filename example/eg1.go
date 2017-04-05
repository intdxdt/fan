package main

import (
	"fmt"
	"time"
	"runtime"
	"fan"
)

const ConCur = 8

var data = []int{40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40}

func Fib(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func streamer() []interface{} {
	var results = make([]interface{}, 0)
	stream := make(chan interface{})
	exit := make(chan struct{})
	defer close(exit)

	//go func() {
	//	time.Sleep(1 * time.Second)
	//	exit <- struct{}{}
	//}()

	go func() {
		for _, d := range data {
			stream <- d
		}
		close(stream)
	}()

	worker := func(v interface{}) interface{} {
		return Fib(v.(int))
	}
	out := fan.IOStream(stream, worker, ConCur, exit)
	for o := range out {
		results = append(results, o.(int))
	}
	return results
}

func streamer_task() []interface{} {
	var results = make([]interface{}, 0)
	stream := make(chan interface{})
	exit := make(chan struct{})
	defer close(exit)

	go func() {
		for _, d := range data {
			stream <- d
		}
		close(stream)
	}()

	worker := func(v interface{}) interface{} {
		return Fib(v.(int))
	}
	_, out := fan.IOTaskStreamer(stream, worker)

	go func() {
		time.Sleep(1 * time.Second)
		//task.Stop()
	}()

	for o := range out {
		results = append(results, o.(int))
	}
	return results
}

func streamer_task_pool() []interface{} {
	var results = make([]interface{}, 0)
	stream := make(chan interface{})
	exit := make(chan struct{})
	defer close(exit)

	//go func() {
	//	time.Sleep(1 * time.Second)
	//	exit <- struct{}{}
	//}()

	go func() {
		for _, d := range data {
			stream <- d
		}
		close(stream)
	}()

	worker := func(v interface{}) interface{} {
		return Fib(v.(int))
	}

	out := fan.IOTaskPool(worker, stream, ConCur, exit)

	for o := range out {
		results = append(results, o.(int))
	}
	return results
}

func pool_payload() []interface{} {
	dat := make([]interface{}, 0)
	exit := make(chan struct{})
	defer  close(exit)

	go func() {
		runtime.Gosched()
		time.Sleep(1 * time.Second)
		exit <- struct{}{}
	}()

	for _, d := range data {
		dat = append(dat, d)
	}

	worker := func(v interface{}) interface{} {
		return Fib(v.(int))
	}
	res := fan.IOTaskPoolPayload(dat, worker, ConCur, exit)
	return res
}

func payload() []interface{} {
	dat := make([]interface{}, 0)
	exit := make(chan struct{})
	defer close(exit)

	//go func() {
	//	runtime.Gosched()
	//	time.Sleep(1 * time.Second)
	//	exit <- struct{}{}
	//}()

	for _, d := range data {
		dat = append(dat, d)
	}

	worker := func(v interface{}) interface{} {
		return Fib(v.(int))
	}

	res := fan.IOPayload(dat, worker, ConCur, exit)
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
	fmt.Println("a -> go ruts:", runtime.NumGoroutine())
	strm := timeIt(streamer, "io streamer")
	fmt.Println(strm)

	fmt.Println("a -> go ruts:", runtime.NumGoroutine())
	strm2 := timeIt(streamer_task, "io streamer task")
	fmt.Println(strm2)

	fmt.Println("a -> go ruts:", runtime.NumGoroutine())
	strm3 := timeIt(streamer_task_pool, "io streamer pool")
	fmt.Println(strm3)

	pool_pyld := timeIt(pool_payload, "TASK pool payload :>")
	fmt.Println(pool_pyld)
	fmt.Println("b -> go ruts:", runtime.NumGoroutine())

	pyld := timeIt(payload, "IO Payload")
	fmt.Println(pyld)
	fmt.Println("b -> go ruts:", runtime.NumGoroutine())

	fmt.Println("exit -> #go :", runtime.NumGoroutine())
	time.Sleep(3 * time.Second)
	done := make(chan struct{})

	time.AfterFunc(1*time.Second, func() {
		fmt.Println("exit -> #go :", runtime.NumGoroutine())
		done <- struct{}{}
		return
	})
	<-done
	fmt.Println("last exit -> #go :", runtime.NumGoroutine())
}
