package main

import (
	"fmt"
	"time"
	"errors"
	"github.com/matryer/runner"
)

var FibData = []int{40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40}

func Fibocci(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return Fibocci(n-1) + Fibocci(n-2)
}

func run_fib_and_stop(in <-chan int, out chan<- int) *runner.Task {
	return runner.Go(func(onStop runner.S) error {
		for val := range in {
			if onStop() {
				break
			}
			out <- Fibocci(val)
		}
		return nil
	})
}

func run_and_stop() {
	var ticker []time.Time
	task := runner.Go(func(shouldStop runner.S) error {
		for {
			ticker = append(ticker, time.Now())
			time.Sleep(100 * time.Millisecond)

			if shouldStop() {
				break
			}
		}
		return nil
	})

	if task.Running() {
		fmt.Println("task is alife")
	}
	time.Sleep(1 * time.Second)

	task.Stop()
	select {
	case <-task.StopChan():
	case <-time.After(2 * time.Second):
		fmt.Println("timed out")
	}
	if !task.Running() {
		fmt.Println("task stoped")
	}
	fmt.Println(len(ticker))
}

func run_err() {

	err := errors.New("something went wrong")
	task := runner.Go(func(shouldStop runner.S) error {
		return err
	})

	time.Sleep(100 * time.Millisecond)
	fmt.Println("task running: is false = ", task.Running())
	fmt.Println(err, task.Err())

	task.Stop()
	select {
	case <-task.StopChan():
	case <-time.After(2 * time.Second):
		fmt.Println("timed out")
	}

}

func main() {
	inbound := make(chan int)
	outbound := make(chan int)
	defer close(outbound)

	task := run_fib_and_stop(inbound, outbound)
	go func() {
		defer close(inbound)
		for _, d := range FibData {
			inbound <- d
			time.Sleep(1 * time.Second)
		}
	}()

	time.AfterFunc(2*time.Second, func() {
		task.Stop()
	})
loop:
	for {
		select {
		case <-task.StopChan():
			fmt.Println("done processing data")
			break loop
		case v := <-outbound:
			fmt.Println("fib val :", v)
		}
	}
	//run_and_stop()
	//run_err()
}
