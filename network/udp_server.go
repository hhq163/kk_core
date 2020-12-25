package network

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/hhq163/kk_core/base"
	"github.com/hhq163/kk_core/common"
)

// type UdpConnSet map[ConnKey]struct{}

type UDPServer struct {
	Addr      string
	NewSess   func(*UDPVconn) Sess
	MaxMsgLen uint32
	UDPConn   *net.UDPConn

	mutexConns sync.Mutex
	wgLn       sync.WaitGroup
	wgConns    sync.WaitGroup

	udpVconns map[uint32]*UDPVconn
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

	u.UDPConn = udpConn
	u.udpVconns = make(map[uint32]*UDPVconn)
}

func (u *UDPServer) run() {
	u.wgLn.Add(1)
	defer u.wgLn.Done()

	recvBuff := make([]byte, u.MaxMsgLen)

	for {
		n, remoteAddr, err := u.UDPConn.ReadFromUDP(recvBuff)
		if err != nil {
			base.Log.Error("UDPConn.ReadFromUDP err=", err.Error())
			break
		}

		if n > 0 {
			packet, err := u.parseHeader(recvBuff[:n])
			if err != nil {
				base.Log.Error("packet parseHeader err=", err.Error())
				continue
			}
			UDPVconn := u.getVTcpConn(remoteAddr, packet.GetUserId())

			UDPVconn.QueuePacket(packet)
		}
	}
}

//消息头解析
func (u *UDPServer) parseHeader(b []byte) (p *common.UdpPacket, err error) {
	if len(b) < 6 {
		return nil, errors.New("package is to short")
	}
	cmd := binary.LittleEndian.Uint32(b[:4])
	msgLen := int(binary.LittleEndian.Uint32(b[4:8]))
	uid := binary.LittleEndian.Uint32(b[8:12])

	if msgLen > len(b) {
		log.Print("message too long from " + strconv.Itoa(int(uid)))
		return nil, errors.New("message too long from " + strconv.Itoa(int(uid)))
	}
	packet := &common.UdpPacket{}
	packet.Init(cmd, uid)
	packet.WriteBytes(b[4:])

	return packet, nil
}

func (u *UDPServer) Close() {
	if u.UDPConn != nil {
		u.UDPConn.Close()
	}
	u.wgLn.Wait()

	u.udpVconns = nil
	u.wgConns.Wait()
}

func (u *UDPServer) getVTcpConn(addr *net.UDPAddr, uid uint32) *UDPVconn {
	var udpVconn *UDPVconn
	if _, ok := u.udpVconns[uid]; ok {
		udpVconn = u.udpVconns[uid]
		udpVconn.UpdateRemote(addr)
	} else {
		udpVconn = newUDPVconn(uid, u.MaxMsgLen, u.UDPConn, addr)
		sess := u.NewSess(udpVconn)
		udpVconn.csession = sess

		u.udpVconns[uid] = udpVconn
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
