package fan

//input source
func src(inputStream []interface{}, exit <-chan struct{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for i := range inputStream {
			select {
			case <- exit :
				return
			case out <- inputStream[i]:
			}
		}
	}()
	return out
}
