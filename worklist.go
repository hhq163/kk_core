package svr_core

import (
	"svr_core/util"
	"sync"
)

type WorkList struct {
	works *util.SyncQueue
	pool  *util.WorkPool
	wg    sync.WaitGroup
}

func NewWorkList(maxGoroutines int) *WorkList {
	w := &WorkList{
		works: util.NewSyncQueue(),
	}
	if maxGoroutines > 0 {
		w.pool = util.NewWorkPool(maxGoroutines)
		w.wg.Add(1)
		go w.Proc()
	}

	return w
}

func (w *WorkList) Push(f func()) {
	w.works.Push(f)
}

//SyncProc proc all work
func (w *WorkList) SyncProc() int {
	fs, _ := w.works.TryPopAll()
	for _, f := range fs {
		f.(func())()
	}
	return len(fs)
}

func (w *WorkList) Proc() {
	defer w.wg.Done()
	for {
		f := w.works.Pop()
		if f == nil {
			return
		}
		w.pool.Run(f.(func()))
	}
}

func (w *WorkList) Close() {
	w.works.Close()
	w.wg.Wait()
	if w.pool != nil {
		w.pool.Shutdown()
	}
}
