package fan

import (
	"github.com/intdxdt/random"
)

type Pool struct {
	workers  map[string]*Worker
	in       <-chan interface{}
	out      chan interface{}
	concur   int
	freechan chan string
	exit     <-chan struct{}
}

//input source
func (p *Pool) Start() <-chan interface{} {
	p.startAllWorkers()
	go func() {
		defer p.stopAllWorkers() //stop all workers on exit
		for {
			select {
			case <-p.exit: //exit the pool
				return
			case o := <-p.in:
				if o == nil { //in closed
					return
				}
				//dispatch work
				p.workers[<-p.freechan].input <- o
			}
		}
	}()
	return p.out
}

func (p *Pool) startAllWorkers() {
	for _, w := range p.workers {
		go w.start()
	}
}

func (p *Pool) stopAllWorkers() {
	for _, w := range p.workers {
		w.stop()
	}
}

func (p *Pool) IsIdle() bool {
	var idle = true
	for _, w := range p.workers {
		idle = idle && w.freeState
		if !idle {
			break
		}
	}
	return idle
}

func NewPool(in <-chan interface{}, fn func(interface{}) interface{}, concur int, exit <-chan struct{}) *Pool {
	var pool = &Pool{
		workers:  make(map[string]*Worker),
		in:       in,
		out:      make(chan interface{}, concur),
		concur:   concur,
		freechan: make(chan string, concur),
		exit:     exit,
	}

	for i := 0; i < concur; i++ {
		w := &Worker{
			id:        random.String(8),
			free:      pool.freechan,
			freeState: false,
			fn:        fn,
			exit:      make(chan struct{}),
			input:     make(chan interface{}),
			out:       pool.out,
		}
		pool.workers[w.id] = w
	}

	return pool
}
