package util_test

import (
	"mangos/util"
	"testing"
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
	qu.SyncProc()
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
