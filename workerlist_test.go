package kk_core

import (
	"testing"

	"github.com/hhq163/kk_core/util"
)

func Benchmark_worklist(b *testing.B) {
	b.StopTimer()
	pool := util.NewWorkerPool(1000, nil)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		pool.Run(func(interface{}) {})
	}
}

func Benchmark_queue(b *testing.B) {
	b.StopTimer()
	qu := NewWorkerList(0, nil)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		qu.Push(func(interface{}) {})
	}
	qu.SyncProcess()
}

func Benchmark_queuelist(b *testing.B) {
	b.StopTimer()
	qu := NewWorkerList(1000, nil)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		qu.Push(func(interface{}) {})
		//qu.Proc()
	}
}
