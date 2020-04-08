package util

import "sync"

type Task func()

type WorkerPool struct {
	work chan Task
	wg   sync.WaitGroup
}

func NewWorkerPool(maxGoroutines int) *WorkerPool {
	p := WorkerPool{
		work: make(chan Task),
	}
	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.work {
				w()
			}
			p.wg.Done()
		}()
	}

	return &p
}

func (p *WorkerPool) Run(w Task) {
	p.work <- w
}
func (p *WorkerPool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
