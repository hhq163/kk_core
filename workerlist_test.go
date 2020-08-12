package kk_core

import (
	"testing"

	"github.com/kk_core/util"
)

func Benchmark_worklist(b *testing.B) {
	b.StopTimer()
	pool := util.NewWorkPool(1000)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		pool.Run(func() {})
	}
}

func Benchmark_queue(b *testing.B) {
	b.StopTimer()
	qu := util.NewWorkList(0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		qu.Push(func() {})
	}
	qu.SyncProcess()
}

func Benchmark_queuelist(b *testing.B) {
	b.StopTimer()
	qu := util.NewWorkList(1000)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		qu.Push(func() {})
		//qu.Proc()
	}
}
