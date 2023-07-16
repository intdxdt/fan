package fan

// Stream takes an in-bound readable stream and
// returns an outbound readable stream
func Stream[T any](stream <-chan T, worker func(T) T,
	concur int, exit <-chan struct{}) <-chan T {
	return process(worker, stream, concur, exit)
}

// Payload performs a task and returns results as a slice of interface{}
func Payload[T any](data []T, worker func(T) T, concur int, exit <-chan struct{}) []T {
	//stage 0 - declare output
	var results = make([]T, 0)

	// stage 1 - input stream
	var in = src(data, exit)

	// stage 2 - process payload
	var res = process(worker, in, concur, exit)

	// stage 3 - store uploads
	var done = store(res, &results, exit)

	<-done //wait

	return results
}
