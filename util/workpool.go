package util

import (
	"sync"
)

type Worker func()

type WorkPool struct {
	work chan Worker
	wg   sync.WaitGroup
}

func NewWorkPool(maxGoroutines int) *WorkPool {
	p := WorkPool{
		work: make(chan Worker),
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
func (p *WorkPool) Run(w Worker) {
	p.work <- w
}
func (p *WorkPool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
