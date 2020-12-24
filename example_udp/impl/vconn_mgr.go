package impl

import (
	"udp_server/base"

	"github.com/hhq163/kk_core/network"
)

type VConnMgr struct {
	udpServer *network.UDPServer
}

func (mgr *VConnMgr) Init() {

	if base.Cfg.UdpAddr != "" {
		mgr.udpServer = new(network.UDPServer)
		mgr.udpServer.Addr = base.Cfg.UdpAddr
		mgr.udpServer.MaxMsgLen = base.Cfg.MaxMsgLen
		mgr.udpServer.NewSess = func(udpVconn *network.UdpVconn) network.Sess {
			u := NewCSession(udpVconn)
			return u
		}

	}

	if mgr.udpServer != nil {
		mgr.udpServer.Start()
	}
}

func (mgr *VConnMgr) Destroy() {
	if mgr.udpServer != nil {
		mgr.udpServer.Close()
	}
}
