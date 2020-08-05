package network

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hhq163/svr_core/base"
)

type WSServer struct {
	Addr        string
	MaxConnNum  int
	MaxMsgLen   uint32
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string
	NewAgent    func(*WSConn) Agent
	ln          net.Listener
	handler     *WSHandler
}

type WSHandler struct {
	maxConnNum int
	maxMsgLen  uint32
	newAgent   func(*WSConn) Agent
	upgrader   websocket.Upgrader
	conns      WebsocketConnSet
	mutexConns sync.Mutex
	wg         sync.WaitGroup
}

func (handler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	protocol := r.Header.Get("Sec-WebSocket-Protocol")
	var header http.Header
	if protocol != "" {
		header = make(http.Header)
		header["Sec-WebSocket-Protocol"] = []string{protocol}
	}
	conn, err := handler.upgrader.Upgrade(w, r, header)
	if err != nil {
		base.Log.Warn("upgrade error: ", err)
		return
	}
	conn.SetReadLimit(int64(handler.maxMsgLen))

	handler.wg.Add(1)
	defer handler.wg.Done()

	handler.mutexConns.Lock()
	if handler.conns == nil {
		handler.mutexConns.Unlock()
		conn.Close()
		return
	}
	if len(handler.conns) >= handler.maxConnNum {
		handler.mutexConns.Unlock()
		conn.Close()
		base.Log.Warn("too many connections")
		return
	}
	handler.conns[conn] = struct{}{}
	handler.mutexConns.Unlock()

	wsConn := newWSConn(conn, handler.maxMsgLen)
	agent := handler.newAgent(wsConn)
	agent.Run()
	// cleanup
	wsConn.Close()
	handler.mutexConns.Lock()
	delete(handler.conns, conn)
	handler.mutexConns.Unlock()
	agent.OnClose()
}

func (server *WSServer) Start() {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		base.Log.Fatal(err)
	}

	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 10000
		base.Log.Info("invalid MaxConnNum, reset to ", server.MaxConnNum)
	}
	if server.MaxMsgLen <= 0 {
		server.MaxMsgLen = 4096
		base.Log.Info("invalid MaxMsgLen, reset to ", server.MaxMsgLen)
	}
	if server.HTTPTimeout <= 0 {
		server.HTTPTimeout = 10 * time.Second
		base.Log.Info("invalid HTTPTimeout, reset to ", server.HTTPTimeout)
	}
	if server.NewAgent == nil {
		base.Log.Fatal("NewAgent must not be nil")
	}
	if server.CertFile != "" || server.KeyFile != "" {
		config := &tls.Config{}
		config.NextProtos = []string{"http/1.1"}

		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(server.CertFile, server.KeyFile)
		if err != nil {
			base.Log.Fatal(err)
		}

		ln = tls.NewListener(ln, config)
	}

	server.ln = ln
	server.handler = &WSHandler{
		maxConnNum: server.MaxConnNum,
		maxMsgLen:  server.MaxMsgLen,
		newAgent:   server.NewAgent,
		conns:      make(WebsocketConnSet),
		upgrader: websocket.Upgrader{
			HandshakeTimeout: server.HTTPTimeout,
			CheckOrigin:      func(_ *http.Request) bool { return true },
			//Subprotocols:     []string{"binary"},
		},
	}

	httpServer := &http.Server{
		Addr:           server.Addr,
		Handler:        server.handler,
		ReadTimeout:    server.HTTPTimeout,
		WriteTimeout:   server.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(ln)
}

func (server *WSServer) Close() {
	server.ln.Close()

	server.handler.mutexConns.Lock()
	for conn := range server.handler.conns {
		conn.Close()
	}
	server.handler.conns = nil
	server.handler.mutexConns.Unlock()

	server.handler.wg.Wait()
}
