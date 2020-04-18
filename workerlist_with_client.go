package svr_core

import (
	"context"
	"sync"

	"github.com/go-redis/redis"
	"github.com/hhq163/svr_core/util"
	"gopkg.in/mgo.v2"
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

func (w *WorkerListWithClient) SyncProc(c interface{}) int {
	fs, _ := w.works.TryPopAll()

	if client, ok := c.(*redis.ClusterClient); ok {
		cClient := client.WithContext(context.Background())
		for _, f := range fs {
			f.(func(c interface{}))(cClient)
		}
	} else if client, ok := c.(*redis.Client); ok {
		cClient := client.WithContext(context.Background())
		for _, f := range fs {
			f.(func(c interface{}))(cClient)
		}
	} else if client, ok := c.(*mgo.Session); ok {
		cClient := client.Clone()
		for _, f := range fs {
			f.(func(c interface{}))(cClient)
		}
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
