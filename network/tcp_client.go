package network

import (
	"net"
	"time"

	"github.com/hhq163/kk_core/base"
	"github.com/hhq163/kk_core/network"
)

type TCPClient struct {
	Addr            string
	ConnectInterval time.Duration
	AutoReconnect   bool
	tcpConn         *TCPConn
	closeFlag       bool
}

func NewTCPClient(addr string) {
	client := new(network.TCPClient)
	client.Addr = addr
	client.AutoReconnect = true
	client.ConnectInterval = 3 * time.Second
}

func (client *TCPClient) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", client.Addr)
		if err == nil || client.closeFlag {
			return conn
		}

		base.Log.Info("connect to ", client.Addr, " error: ", err)
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *TCPClient) connect() {

reconnect:
	conn := client.dial()
	if conn == nil {
		base.Log.Info("client.dial return nil")
		return
	}
	client.tcpConn = NewTCPConn(conn, client.MaxMsgLen)

	tcpConn := NewTCPConn(conn, client.MaxMsgLen)
	agent := client.NewAgent(tcpConn)
	client.Run()

	if client.AutoReconnect {
		time.Sleep(client.ConnectInterval)
		goto reconnect
	}
}

func (client *TCPClient) Run() {
	go func() {
		// PING := &common.WorldPacket{}
		// PING.Initialize(CMSG_KEEP_ALIVE)
		// for {
		// 	if !a.conn.IsClosed() {
		// 		a.conn.WriteMsg(PING)
		// 	}
		// 	time.Sleep(1 * time.Second)
		// }
	}()

	for {
		packet, err := a.conn.ReadMsg()
		if err != nil {
			slog.Error("register read err=", err)
			return
		}
		// if packet.GetOpCode() < MSG_NULL_ACTION || packet.GetOpCode() > NUM_MSG_TYPES {
		// 	slog.Info("Signal recv unknow opcode:", packet.GetOpCode())
		// 	continue
		// }
		// sGameManager.QueuePacket(packet)
	}
}

func (client *TCPClient) Close() {
	client.closeFlag = true
	client.tcpConn.Close()
	client.tcpConn = nil
}

func (client *TCPClient) IsClosed() bool {
	return client.closeFlag
}
