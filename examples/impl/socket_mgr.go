package impl

import (
	"kk_server/base"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hhq163/kk_core/network"
)

type SocketMgr struct {
	wsServer *network.WSServer
}

func (mgr *SocketMgr) Init() {

	if base.Cfg.WSAddr != "" {
		mgr.wsServer = new(network.WSServer)
		mgr.wsServer.Addr = base.Cfg.WSAddr
		mgr.wsServer.MaxConnNum = base.Cfg.MaxConnNum
		mgr.wsServer.MaxMsgLen = base.Cfg.MaxMsgLen
		mgr.wsServer.HTTPTimeout = base.Cfg.HttpTimeout * time.Second
		mgr.wsServer.CertFile = base.Cfg.CertFile
		mgr.wsServer.KeyFile = base.Cfg.KeyFile
		mgr.wsServer.NewAgent = func(conn net.Conn) network.Agent {
			tcpConn := network.NewTCPConn(conn, base.Cfg.MaxMsgLen)
			a := &agent{conn: tcpConn}
			return a
		}
		mgr.wsServer.Handler = &WSHandler{
			maxConnNum: base.Cfg.MaxConnNum,
			maxMsgLen:  base.Cfg.MaxMsgLen,
			newAgent:   mgr.wsServer.NewAgent,
			conns:      make(network.WsConnSet),
			upgrader: websocket.Upgrader{
				HandshakeTimeout: base.Cfg.HttpTimeout * time.Second,
				CheckOrigin:      func(_ *http.Request) bool { return true },
				//Subprotocols:     []string{"binary"},
			},
		}
	}

	if mgr.wsServer != nil {
		mgr.wsServer.Start()
	}
}

func (mgr *SocketMgr) Destroy() {
	if mgr.wsServer != nil {
		mgr.wsServer.Close()
	}
}
