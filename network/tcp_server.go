package network

import (
	"net"
	"sync"

	"github.com/hhq163/kk_core/base"
)

//ConnSet server conn map
type ConnSet map[net.Conn]struct{}

//TCPServer TCP server manager
type TCPServer struct {
	Addr       string
	MaxConnNum int
	NewAgent   func(*TCPConn) Agent
	MaxMsgLen  uint32
	ln         net.Listener
	conns      ConnSet
	mutexConns sync.Mutex
	wgLn       sync.WaitGroup
	wgConns    sync.WaitGroup
}

//Start start server
func (server *TCPServer) Start() {
	server.Init()
	go server.run()
}

func (server *TCPServer) Init() {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		base.Log.Fatal(err)
		return
	}

	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 10000
		base.Log.Info("invalid MaxConnNum, reset to ", server.MaxConnNum)
		return
	}
	if server.NewAgent == nil {
		base.Log.Fatal("NewAgent must not be nil")
		return
	}

	server.ln = ln
	server.conns = make(ConnSet)
}

func (server *TCPServer) run() {
	server.wgLn.Add(1)
	defer server.wgLn.Done()

	for {
		conn, err := server.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				base.Log.Info("accept temp err=", ne.Error())
				continue
			}
			return
		}

		server.mutexConns.Lock()
		if len(server.conns) >= server.MaxConnNum {
			server.mutexConns.Unlock()
			conn.Close()
			base.Log.Warn("too many connections")
			continue
		}
		server.conns[conn] = struct{}{}
		server.mutexConns.Unlock()

		server.wgConns.Add(1)

		tcpConn := newTCPConn(conn, server.MaxMsgLen)
		agent := server.NewAgent(tcpConn)
		go func() {
			agent.Run()

			// cleanup
			tcpConn.Close()
			server.mutexConns.Lock()
			delete(server.conns, conn)
			server.mutexConns.Unlock()
			agent.OnClose()
			server.wgConns.Done()
		}()
	}
}

func (server *TCPServer) Close() {
	server.ln.Close()
	server.wgLn.Wait()

	server.mutexConns.Lock()
	for conn := range server.conns {
		conn.Close()
	}
	server.conns = nil
	server.mutexConns.Unlock()
	server.wgConns.Wait()
}
