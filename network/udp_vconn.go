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
type UDPVconn struct {
	uid        uint32
	remote     *net.UDPAddr
	UDPConn    *net.UDPConn
	writeQueue *util.SyncQueue
	csession   Sess

	maxMsgLen  uint32
	closeFlag  int32
	activeTime int64 //the time of last receive msg
	Crypt      auth.AuthCrypt
}

func newUDPVconn(uid uint32, maxMsgLen uint32, udpConn *net.UDPConn, remote *net.UDPAddr) *UDPVconn {
	udpVconn := &UDPVconn{
		uid:        uid,
		remote:     remote,
		UDPConn:    udpConn,
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

func (u *UDPVconn) IsClosed() bool {
	return atomic.LoadInt32(&u.closeFlag) == 1
}

func (u *UDPVconn) Close() {
	u.writeQueue.Close()
}

func (u *UDPVconn) Write(b []byte) {
	u.writeQueue.Push(b)
}

func (u *UDPVconn) ReadMsg(b []byte) (err error) {
	cmd := binary.LittleEndian.Uint16(b[:2])
	msgLen := int(binary.LittleEndian.Uint16(b[2:4]))

	if msgLen > len(b) {
		log.Print("message too long from " + u.remote.String())
		return errors.New("message too long from " + u.remote.String())
	}
	packet := &common.Packet{}
	packet.Initialize(cmd)
	packet.WriteBytes(b[4:])

	return nil
}

func (u *UDPVconn) QueuePacket(msg common.IPacket) {
	u.csession.QueuePacket(msg)
}

func (u *UDPVconn) DirectWrite(b []byte) {
	u.UDPConn.Write(b)
}

func (u *UDPVconn) WriteMsg(packet common.IPacket) error {
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

func (u *UDPVconn) IsTimeout(maxTime uint32) bool {
	var ret bool
	cur := uint32(time.Now().Unix() - u.activeTime)
	if cur > maxTime {
		ret = true
	}
	return ret
}
