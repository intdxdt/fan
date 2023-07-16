package fan

func src[T any](data []T, exit <-chan struct{}) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for i := range data {
			select {
			case <-exit:
				return
			case out <- data[i]:
			}
		}
	}()
	return out
}
