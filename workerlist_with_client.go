package svr_core

import (
	"sync"

	"github.com/hhq163/svr_core/util"
)

type WorkerListWithClient struct {
	works *util.SyncQueue
	pool  *util.WorkerPoolWithClient
	wg    sync.WaitGroup
}

func NewWorkerListWithClient(maxGoroutines int, c interface{}) *WorkerListWithClient {
	w := &WorkerListWithClient{
		works: util.NewSyncQueue(),
	}
	if maxGoroutines > 0 {
		w.pool = util.NewWorkerPoolWithClient(maxGoroutines, c)
		w.wg.Add(1)
		go w.Proc()
	}

	return w
}

func (w *WorkerListWithClient) Push(f func(i interface{})) {
	w.works.Push(f)
}

func (w *WorkerListWithClient) SyncProc() int {
	fs, _ := w.works.TryPopAll()
	for _, f := range fs {
		f.(func(c interface{}))(false)
	}
	return len(fs)
}

func (w *WorkerListWithClient) Proc() {
	defer w.wg.Done()
	for {
		f := w.works.Pop()
		if f == nil {
			return
		}
		w.pool.Run(f.(func(c interface{})))
	}
}

func (w *WorkerListWithClient) Close() {
	w.works.Close()
	w.wg.Wait()
	if w.pool != nil {
		w.pool.Shutdown()
	}
}
