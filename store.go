package fan

//store
func store(in <-chan interface{}, results *[]interface{}, exit <-chan struct{}) <-chan struct{} {
	//create done signal
	done := make(chan struct{})
	//start new go routine to listen for results
	go func() {
		defer close(done) //close when res chan closes
		for data := range in {
			select {
			//listen to results channel
			case <-exit:
				return
			default:
				*results = append(*results, data)
			}
		}
	}()
	//return the wait signal to whoever wants to know when all is done
	return done
}
