package util

import "sync"

//WorkList work list
type WorkList struct {
	works *SyncQueue
	pool  *WorkPool
	wg    sync.WaitGroup
}

//NewWorkList new WorList
func NewWorkList(maxGoroutines int) *WorkList {
	w := &WorkList{
		works: NewSyncQueue(),
	}
	if maxGoroutines > 0 {
		w.pool = NewWorkPool(maxGoroutines)
		w.wg.Add(1)
		go w.Proc()
	}

	return w
}

//Push push a work to list
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

//Proc proc work
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

//Close close queue
func (w *WorkList) Close() {
	w.works.Close()
	w.wg.Wait()
	if w.pool != nil {
		w.pool.Shutdown()
	}
}
