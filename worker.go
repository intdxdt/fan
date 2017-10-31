package fan

type Worker struct {
	id        string
	free      chan<- string
	freeState bool
	fn        func(interface{}) interface{}
	exit      chan struct{}
	input     chan interface{}
	out       chan<- interface{}
	running   bool
}

func (w *Worker) isBusy() {
	w.freeState = false
}

func (w *Worker) isFree() {
	w.freeState = true
	w.free <- w.id
}

func (w *Worker) start() {
	w.running = true
	w.isFree() //first: announce you are free
	defer func() {
		w.isFree()
		w.running = false
	}()
	for {
		select {
		case <-w.exit:
			return
		case o := <-w.input:
			w.isBusy()       //busy
			w.out <- w.fn(o) //do the work
			w.isFree()       //free
		}
	}
}

func (w *Worker) stop() {
	close(w.exit)
}
