package fan

type Worker[T any] struct {
	id        string
	free      chan<- string
	freeState bool
	fn        func(T) T
	exit      chan struct{}
	input     chan T
	out       chan<- T
	running   bool
}

func (w *Worker[T]) isBusy() {
	w.freeState = false
}

func (w *Worker[T]) isFree() {
	w.freeState = true
	w.free <- w.id
}

func (w *Worker[T]) start() {
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

func (w *Worker[T]) stop() {
	close(w.exit)
}
