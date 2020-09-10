package impl

import (
	"github.com/hhq163/kk_core"
)

var WorldList = kk_core.NewWorkerList(0)
var mysqlworklist *kk_core.WorkerList

func InitWorkers() {
	mysqlworklist = kk_core.NewWorkerList(20)
}
func StopWorkers() {
	WorldList.Close()
}

func PushMysql(f func()) {
	mysqlworklist.Push(f)
}

func PushWorld(f func()) {
	WorldList.Push(f)
}
