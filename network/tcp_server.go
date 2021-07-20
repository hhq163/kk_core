package network

import (
	"net"
	"sync"
	"time"

	"github.com/hhq163/kk_core/base"
)

//ConnSet server conn map
type ConnSet map[net.Conn]struct{}

//TCPServer TCP server manager
type TCPServer struct {
	Addr       string
	MaxConnNum int
	NewAgent   func(net.Conn) Agent
	MaxMsgLen  uint32
	ln         net.Listener
	conns      ConnSet
	mutexConns sync.Mutex
	wgLn       sync.WaitGroup
	wgConns    sync.WaitGroup
}

//Start start server
func (server *TCPServer) Start(maxConnNum int) {
	server.Init(maxConnNum)
	go server.run()
}

func (server *TCPServer) Init(maxConnNum int) {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		base.Log.Fatal(err)
		return
	}
	if maxConnNum == 0 {
		maxConnNum = 10000
	}

	if server.MaxConnNum <= 0 {
		server.MaxConnNum = maxConnNum
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

	var tempDelay time.Duration
	for {
		conn, err := server.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				base.Log.Info("accept error: ", err, "; retrying in ", tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0

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

		agent := server.NewAgent(conn)
		go func() {
			agent.Run()

			// cleanup
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
