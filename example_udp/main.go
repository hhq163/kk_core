package main

import (
	"log"
	"net/http"

	"udp_server/impl"

	"github.com/hhq163/kk_core/base"
)

func main() {
	impl.InitWorkers()

	base.LogInit(true, "udp_server")

	base.Log.Info("server started")
	socketmgr := new(impl.VConnMgr)
	socketmgr.Init()
	impl.Start()

}

func profile() {
	log.Fatal(http.ListenAndServe("0.0.0.0:5001", nil))
}
