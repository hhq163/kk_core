package svr_core

import (
	"sync"

	"github.com/hhq163/svr_core/util"
)

//根据uid进行投递任务，同一用户任务顺序执行
type WorkerGroup struct {
	WorkGroup []*WorkerListWithClient
	Length    int
}

func NewWorkGroup(groupNum int, c interface{}) *WorkerGroup {
	wg := &WorkerGroup{
		WorkGroup: make([]*WorkerListWithClient, groupNum),
		Length:    groupNum,
	}
	for k := range wg.WorkGroup {
		wg.WorkGroup[k] = NewWorkerListWithClient(1, c)
	}

	return wg
}

func (wg *WorkerGroup) Push(uid int, f func(i interface{})) {
	index := uid % wg.Length
	wg.WorkGroup[index].Push(f)
}

//在协程池中使用固定db连接
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
		f.(func())()
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
