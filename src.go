package fan

//input source
func src(data []interface{}, exit <-chan struct{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for i := range data {
			select {
			case <- exit :
				return
			case out <- data[i]:
			}
		}
	}()
	return out
}
