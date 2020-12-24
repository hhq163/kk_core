package network

import (
	"encoding/binary"
	"net"
	"sync"

	"github.com/hhq163/kk_core/base"
)

type UdpConnSet map[ConnKey]struct{}

type UDPServer struct {
	Addr      string
	NewSess   func(*UdpVconn) Sess
	MaxMsgLen uint32
	UdpConn   *net.UDPConn

	mutexConns sync.Mutex
	wgLn       sync.WaitGroup
	wgConns    sync.WaitGroup

	udpVconns map[ConnKey]*UdpVconn
}

func (u *UDPServer) Start() {
	u.Init()
	go u.run()
}

func (u *UDPServer) Init() {
	addr, err := net.ResolveUDPAddr("udp", u.Addr)
	if err != nil {
		base.Log.Fatal(err)
		return
	}

	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		base.Log.Fatal(err)
		return
	}

	u.UdpConn = udpConn
	u.udpVconns = make(map[ConnKey]*UdpVconn)
}

func (u *UDPServer) run() {
	u.wgLn.Add(1)
	defer u.wgLn.Done()

	recvBuff := make([]byte, u.MaxMsgLen)

	for {
		n, remoteAddr, err := u.UdpConn.ReadFromUDP(recvBuff)
		if err != nil {
			base.Log.Error("UdpConn.ReadFromUDP err=", err.Error())
			break
		}

		if n > 0 {
			UdpVconn := u.getVTcpConn(remoteAddr)

			UdpVconn.ReadMsg(recvBuff[:n])
		}
	}
}

func (u *UDPServer) Close() {
	if u.UdpConn != nil {
		u.UdpConn.Close()
	}
	u.wgLn.Wait()

	u.udpVconns = nil
	u.wgConns.Wait()
}

func (u *UDPServer) getVTcpConn(addr *net.UDPAddr) *UdpVconn {
	key := genConnKey(addr)
	var udpVconn *UdpVconn
	if _, ok := u.udpVconns[key]; ok {
		udpVconn = u.udpVconns[key]
	} else {
		udpVconn = newUdpVconn(key, u.MaxMsgLen, u.UdpConn, addr)
		sess := u.NewSess(udpVconn)
		udpVconn.csession = sess

		u.udpVconns[key] = udpVconn
	}
	return udpVconn
}

//A net.Addr where the IP is split into two fields so you can use it as a key
type ConnKey struct {
	IPHigh uint64
	IPLow  uint64
	Port   int
}

func genConnKey(addr *net.UDPAddr) ConnKey {
	if len(addr.IP) == net.IPv4len {
		return ConnKey{
			IPHigh: 0,
			IPLow:  uint64(binary.BigEndian.Uint32(addr.IP)),
			Port:   addr.Port,
		}
	}

	return ConnKey{
		IPHigh: binary.BigEndian.Uint64(addr.IP[:8]),
		IPLow:  binary.BigEndian.Uint64(addr.IP[8:]),
		Port:   addr.Port,
	}
}
