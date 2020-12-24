package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/hhq163/kk_core/auth"
	"github.com/hhq163/kk_core/common"
	"github.com/hhq163/kk_core/util"
)

//UDP虚拟连接类
type UdpVconn struct {
	connkey    ConnKey
	remote     *net.UDPAddr
	UdpConn    *net.UDPConn
	writeQueue *util.SyncQueue
	csession   Sess

	maxMsgLen  uint32
	closeFlag  int32
	activeTime int64 //the time of last receive msg
	Crypt      auth.AuthCrypt
}

func newUdpVconn(connkey ConnKey, maxMsgLen uint32, udpConn *net.UDPConn, remote *net.UDPAddr) *UdpVconn {
	udpVconn := &UdpVconn{
		connkey:    connkey,
		remote:     remote,
		UdpConn:    udpConn,
		writeQueue: util.NewSyncQueue(),
		maxMsgLen:  maxMsgLen,
		activeTime: time.Now().Unix(),
	}

	go func() {
		var tmp bytes.Buffer
		for {
			time.Sleep(100 * time.Millisecond)
			bs := udpVconn.writeQueue.PopAll()
			if bs == nil {
				break
			}
			for _, b := range bs {
				buf, ok := b.([]byte)
				if ok {
					tmp.Write(buf)
				}
			}
			if tmp.Len() != 0 {
				_, err := udpConn.WriteToUDP(tmp.Bytes(), remote)
				if err != nil {
					log.Print(err.Error())
					goto closeSocket
				}
			}
			tmp.Reset()
		}
	closeSocket:
		atomic.StoreInt32(&udpVconn.closeFlag, 1)
	}()

	return udpVconn
}

func (u *UdpVconn) IsClosed() bool {
	return atomic.LoadInt32(&u.closeFlag) == 1
}

func (u *UdpVconn) Close() {
	u.writeQueue.Close()
}

func (u *UdpVconn) Write(b []byte) {
	u.writeQueue.Push(b)
}

func (u *UdpVconn) ReadMsg(b []byte) (err error) {
	cmd := binary.LittleEndian.Uint16(b[:2])
	msgLen := int(binary.LittleEndian.Uint16(b[2:4]))

	if msgLen > len(b) {
		log.Print("message too long from " + u.remote.String())
		return errors.New("message too long from " + u.remote.String())
	}
	packet := &common.Packet{}
	packet.Initialize(cmd)
	packet.WriteBytes(b[4:])

	u.csession.QueuePacket(packet)
	return nil
}

func (u *UdpVconn) DirectWrite(b []byte) {
	u.UdpConn.Write(b)
}

func (u *UdpVconn) WriteMsg(packet common.IPacket) error {
	if u.IsClosed() {
		return errors.New("udp conn is closed")
	}
	// get len
	msgLen := uint16(packet.Len() + int(mLen) + 2)
	header := new(bytes.Buffer)
	binary.Write(header, binary.LittleEndian, packet.GetCmd())
	binary.Write(header, binary.LittleEndian, msgLen)
	binary.Write(header, binary.LittleEndian, packet.Bytes())
	u.Write(header.Bytes())
	return nil
}

func (u *UdpVconn) IsTimeout(maxTime uint32) bool {
	var ret bool
	cur := uint32(time.Now().Unix() - u.activeTime)
	if cur > maxTime {
		ret = true
	}
	return ret
}
