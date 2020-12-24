package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync/atomic"
	"time"

	"log"

	"github.com/hhq163/kk_core/auth"
	"github.com/hhq163/kk_core/common"

	"github.com/hhq163/kk_core/util"
)

const mLen uint32 = 4 //长度占用的字节数

//TCPConn tcp连接类
type TCPConn struct {
	conn       net.Conn
	writeQueue *util.SyncQueue
	maxMsgLen  uint32
	closeFlag  int32
	activeTime int64 //the time of last receive msg
	Crypt      auth.AuthCrypt
}

func newTCPConn(conn net.Conn, maxMsgLen uint32) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.conn = conn
	tcpConn.writeQueue = util.NewSyncQueue()
	tcpConn.maxMsgLen = maxMsgLen
	//tcp写goroutine
	go func() {
		var tmp bytes.Buffer
		for {
			time.Sleep(100 * time.Millisecond)
			bs := tcpConn.writeQueue.PopAll()
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
				_, err := conn.Write(tmp.Bytes())
				if err != nil {
					log.Print(err.Error())
					goto closeSocket
				}
			}
			tmp.Reset()
		}
	closeSocket:
		atomic.StoreInt32(&tcpConn.closeFlag, 1)
		conn.Close()
	}()

	return tcpConn
}

func (tcpConn *TCPConn) IsClosed() bool {
	return atomic.LoadInt32(&tcpConn.closeFlag) == 1
}
func (tcpConn *TCPConn) Close() {
	tcpConn.writeQueue.Close()
}

func (tcpConn *TCPConn) Write(b []byte) {
	tcpConn.writeQueue.Push(b)
}

//不进队列，直接发送
func (tcpConn *TCPConn) DirectWrite(b []byte) {
	tcpConn.conn.Write(b)
}
func (tcpConn *TCPConn) Read(b []byte) (n int, err error) {
	return tcpConn.conn.Read(b)
}

func (tcpConn *TCPConn) LocalAddr() net.Addr {
	return tcpConn.conn.LocalAddr()
}

func (tcpConn *TCPConn) RemoteAddr() net.Addr {
	return tcpConn.conn.RemoteAddr()
}

func (tcpConn *TCPConn) ReadMsg() (common.IPacket, error) {
	headbuf := make([]byte, mLen)
	// read len
	if _, err := io.ReadFull(tcpConn, headbuf); err != nil {
		return nil, err
	}
	tcpConn.Crypt.DecryptRecv(headbuf)
	msgLen := uint32(binary.LittleEndian.Uint16(headbuf[:mLen]))
	cmdId := binary.LittleEndian.Uint16(headbuf[mLen : mLen+2])
	// check len
	if msgLen > tcpConn.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < mLen {
		return nil, errors.New("message too short")
	}

	// data
	if msgLen > mLen {
		msgData := make([]byte, msgLen-mLen)
		if _, err := io.ReadFull(tcpConn, msgData); err != nil {
			return nil, err
		}
		packet := &common.Packet{}
		packet.Initialize(cmdId)
		packet.WriteBytes(msgData)
		return packet, nil
	}
	packet := &common.Packet{}
	packet.Initialize(cmdId)
	return packet, nil

}

func (tcpConn *TCPConn) WriteMsg(packet common.IPacket) error {
	if tcpConn.IsClosed() {
		return errors.New("socket is closed")
	}
	// get len
	msgLen := uint16(packet.Len() + int(mLen) + 2)
	header := new(bytes.Buffer)
	binary.Write(header, binary.LittleEndian, msgLen)
	binary.Write(header, binary.LittleEndian, packet.GetCmd())
	tcpConn.Crypt.EncryptSend(header.Bytes())
	binary.Write(header, binary.LittleEndian, packet.Bytes())
	tcpConn.Write(header.Bytes())
	return nil
}

func (tcpConn *TCPConn) InitCrypt(k []byte) {
	tcpConn.Crypt.Init(k)
}

func (tcpConn *TCPConn) IsTimeout(maxTime uint32) bool {
	var ret bool
	cur := uint32(time.Now().Unix() - tcpConn.activeTime)
	if cur > maxTime {
		ret = true
	}
	return ret
}
