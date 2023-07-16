package fan

func store[T any](in <-chan T, results *[]T, exit <-chan struct{}) <-chan struct{} {
	var done = make(chan struct{})
	go func() {
		defer close(done)
		for data := range in {
			select {
			case <-exit:
				return
			default:
				*results = append(*results, data)
			}
		}
	}()
	return done
}
