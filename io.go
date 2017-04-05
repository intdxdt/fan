package fan

import (
	"github.com/matryer/runner"
)

//IO Stream takes an in-bound readable readable stream and
// returns an outbound readable stream
func IOStream(stream <-chan interface{}, worker func(interface{}) interface{},
	concur int, exit <-chan struct{}) <-chan interface{} {
	return process(worker, stream, concur, exit)
}

func IOTaskStreamer(in <-chan interface{}, worker func(interface{}) interface{}) (*runner.Task, <-chan interface{}) {
	return processTask(in, worker)
}

//IOTaskPool taks a worker , input stream and returns a readable output stream
func IOTaskPool(worker func(interface{}) interface{},
	in <-chan interface{}, concurrency int, exit <-chan struct{}) <-chan interface{} {
	return processTaskPool(worker, in, concurrency, exit)
}

//Process data using a pool of tasks
func IOTaskPoolPayload( data []interface{}, worker func(interface{}) interface{}, concur int,
	exit <-chan struct{}) []interface{}  {
	//return processTaskPool(worker, in, concurrency, exit)
		//stage 0 - declare output
	var results = make([]interface{}, 0)

	// stage 1 - input stream
	in := src(data, exit)

	// stage 2 - process payload
	res := processTaskPool(worker, in, concur, exit)

	// stage 3 - store uploads
	done := store(res, &results, exit)

	<-done //wait for done signal
	return results
}

//IO Payload performs a task and returns results as a slice of interface{}
func IOPayload(
	data []interface{}, worker func(interface{}) interface{}, concur int,
	exit <-chan struct{}) []interface{} {
	//stage 0 - declare output
	var results = make([]interface{}, 0)

	// stage 1 - input stream
	in := src(data, exit)

	// stage 2 - process payload
	res := process(worker, in, concur, exit)

	// stage 3 - store uploads
	done := store(res, &results, exit)

	<-done //wait for done signal
	return results
}