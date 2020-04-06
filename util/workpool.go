package util

import (
	"sync"
)

type Worker func()

type WorkPool struct {
	work chan Worker //通道，发送worker
	wg   sync.WaitGroup
}

func NewWorkPool(maxGoroutines int) *WorkPool {
	p := WorkPool{
		work: make(chan Worker),
	}
	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.work { //这种方式从一个chan中取东西， select 是从多个chan中取东西，如果没有则阻塞，只有在通道被关闭时才会终止循环
				w() //执行
			}
			p.wg.Done()
		}() //多协程运行，
	}
	return &p
}
func (p *WorkPool) Run(w Worker) { //向通道发送work，协程运行
	p.work <- w
}
func (p *WorkPool) Shutdown() {
	close(p.work) //关闭通道，协程退出循环
	p.wg.Wait()   //等待协程终止
}
