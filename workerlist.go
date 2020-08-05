package svr_core

import (
	"sync"

	"github.com/hhq163/svr_core/util"
)

type WorkerList struct {
	works *util.SyncQueue
	pool  *util.WorkerPool
	wg    sync.WaitGroup
}

func NewWorkerList(maxGoroutines int) *WorkerList {
	w := &WorkerList{
		works: util.NewSyncQueue(),
	}
	if maxGoroutines > 0 {
		w.pool = util.NewWorkerPool(maxGoroutines)
		w.wg.Add(1)
		go w.Proc()
	}

	return w
}

func (w *WorkerList) Push(f func()) {
	w.works.Push(f)
}

//SyncProc 执行所有任务
func (w *WorkerList) SyncProc() int {
	fs, _ := w.works.TryPopAll()
	for _, f := range fs {
		f.(func())()
	}
	return len(fs)
}

func (w *WorkerList) Proc() {
	defer w.wg.Done()
	for {
		f := w.works.Pop()
		if f == nil {
			return
		}
		w.pool.Run(f.(func()))
	}
}

func (w *WorkerList) Close() {
	w.works.Close()
	w.wg.Wait()
	if w.pool != nil {
		w.pool.Shutdown()
	}
}
