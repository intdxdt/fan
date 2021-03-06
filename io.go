package fan

//IO Stream takes an in-bound readable readable stream and
// returns an outbound readable stream
func Stream(stream <-chan interface{}, worker func(interface{}) interface{},
	concur int, exit <-chan struct{}) <-chan interface{} {
	return process(worker, stream, concur, exit)
}

//IO Payload performs a task and returns results as a slice of interface{}
func Payload(
	data []interface{}, worker func(interface{}) interface{}, concur int,
	exit <-chan struct{}) []interface{} {
	//stage 0 - declare output
	var results = make([]interface{}, 0)

	// stage 1 - input stream
	var in = src(data, exit)

	// stage 2 - process payload
	var res = process(worker, in, concur, exit)

	// stage 3 - store uploads
	var done = store(res, &results, exit)

	<-done //wait

	return results
}
