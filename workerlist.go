package kk_core

import (
	"sync"

	"github.com/hhq163/kk_core/util"
)

//根据uid进行投递任务，同一用户任务顺序执行
type WorkerGroup struct {
	WorkGroup []*WorkerList
	Length    int
}

func NewWorkGroup(groupNum int, c interface{}) *WorkerGroup {
	wg := &WorkerGroup{
		WorkGroup: make([]*WorkerList, groupNum),
		Length:    groupNum,
	}
	for k := range wg.WorkGroup {
		wg.WorkGroup[k] = NewWorkerList(1, c)
	}

	return wg
}

func (wg *WorkerGroup) Push(uid int, f func(i interface{})) {
	index := uid % wg.Length
	wg.WorkGroup[index].Push(f)
}

type WorkerList struct {
	works *util.SyncQueue
	pool  *util.WorkerPool
	wg    sync.WaitGroup
}

func NewWorkerList(maxGoroutines int, c interface{}) *WorkerList {
	w := &WorkerList{
		works: util.NewSyncQueue(),
	}
	if maxGoroutines > 0 {
		w.pool = util.NewWorkerPool(maxGoroutines, c)
		w.wg.Add(1)
		go w.Process()
	}

	return w
}

func (w *WorkerList) Push(f func(i interface{})) {
	w.works.Push(f)
}

//在主协程顺序执行所有任务
func (w *WorkerList) SyncProcess() int {
	fs, _ := w.works.TryPopAll()
	for _, f := range fs {
		f.(func())()
	}
	return len(fs)
}

func (w *WorkerList) Process() {
	defer w.wg.Done()
	for {
		f := w.works.Pop()
		if f == nil {
			return
		}
		w.pool.Run(f.(func(c interface{})))
	}
}

func (w *WorkerList) Close() {
	w.works.Close()
	w.wg.Wait()
	if w.pool != nil {
		w.pool.Shutdown()
	}
}
