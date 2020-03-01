package fan

func store(in <-chan interface{}, results *[]interface{}, exit <-chan struct{}) <-chan struct{} {
	done := make(chan struct{})
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
