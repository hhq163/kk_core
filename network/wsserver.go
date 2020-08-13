package network

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hhq163/kk_core/base"
)

type WsConnSet map[*websocket.Conn]struct{}

type KKHandle interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Close()
}

type WSServer struct {
	Addr        string
	MaxConnNum  int
	MaxMsgLen   uint32
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string
	NewAgent    func(*Conn) Agent
	ln          net.Listener
	Handler     *KKHandle
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

	httpServer := &http.Server{
		Addr:           server.Addr,
		Handler:        *server.Handler,
		ReadTimeout:    server.HTTPTimeout,
		WriteTimeout:   server.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(ln)
}

func (server *WSServer) Close() {
	server.ln.Close()

	server.Handler.Close()
}
