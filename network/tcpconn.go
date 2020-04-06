package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"mangos/common"
	"mangos/core/auth"
	"mangos/core/slog"
	"net"
	"sync/atomic"
	"time"

	"mangos/core/util"
)

//ConnSet server conn map
type ConnSet map[net.Conn]struct{}

const HeaderLen uint32 = 4

//TCPConn tcp连接类
type TCPConn struct {
	conn      net.Conn
	writeChan *util.SyncQueue
	maxMsgLen uint32
	closeFlag int32
	Crypt     auth.AuthCrypt
}

func newTCPConn(conn net.Conn, maxMsgLen uint32) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.conn = conn
	tcpConn.writeChan = util.NewSyncQueue()
	tcpConn.maxMsgLen = maxMsgLen
	//tcp写goroutine
	go func() {
		var tmp bytes.Buffer
		for {
			time.Sleep(100 * time.Millisecond)
			bs := tcpConn.writeChan.PopAll()
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
					slog.Info(err)
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
	tcpConn.writeChan.Close()
}

func (tcpConn *TCPConn) Write(b []byte) {
	tcpConn.writeChan.Push(b)
}
func (tcpConn *TCPConn) DirectWrite(b []byte) {
	tcpConn.conn.Write(b)
}

func (tcpConn *TCPConn) Read(b []byte) (int, error) {
	return tcpConn.conn.Read(b)
}

func (tcpConn *TCPConn) LocalAddr() net.Addr {
	return tcpConn.conn.LocalAddr()
}

func (tcpConn *TCPConn) RemoteAddr() net.Addr {
	return tcpConn.conn.RemoteAddr()
}

func (tcpConn *TCPConn) ReadMsg() (*common.WorldPacket, error) {
	b := make([]byte, 4)
	// read len
	if _, err := io.ReadFull(tcpConn, b); err != nil {
		return nil, err
	}
	tcpConn.Crypt.DecryptRecv(b)
	msgLen := uint32(binary.LittleEndian.Uint16(b[:2]))
	opCode := binary.LittleEndian.Uint16(b[2:])
	// check len
	if msgLen > tcpConn.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < HeaderLen {
		return nil, errors.New("message too short")
	}

	// data
	if msgLen > 4 {
		msgData := make([]byte, msgLen-4)
		if _, err := io.ReadFull(tcpConn, msgData); err != nil {
			return nil, err
		}
		packet := &common.WorldPacket{}
		packet.Initialize(opCode)
		packet.WriteBytes(msgData)
		return packet, nil
	}
	packet := &common.WorldPacket{}
	packet.Initialize(opCode)
	return packet, nil

}

func (tcpConn *TCPConn) WriteMsg(packet *common.WorldPacket) error {
	if tcpConn.IsClosed() {
		return errors.New("socket is closed")
	}
	// get len
	msgLen := uint16(packet.Len() + 4)
	header := new(bytes.Buffer)
	binary.Write(header, binary.LittleEndian, msgLen)
	binary.Write(header, binary.LittleEndian, packet.GetOpCode())
	tcpConn.Crypt.EncryptSend(header.Bytes())
	binary.Write(header, binary.LittleEndian, packet.Bytes())
	tcpConn.Write(header.Bytes())
	return nil
}
func (tcpConn *TCPConn) DirectWriteMsg(packet *common.WorldPacket) error {
	if tcpConn.IsClosed() {
		return errors.New("socket is closed")
	}
	// get len
	msgLen := uint16(packet.Len() + 4)
	header := new(bytes.Buffer)
	binary.Write(header, binary.LittleEndian, msgLen)
	binary.Write(header, binary.LittleEndian, packet.GetOpCode())
	tcpConn.Crypt.EncryptSend(header.Bytes())
	binary.Write(header, binary.LittleEndian, packet.Bytes())
	tcpConn.DirectWrite(header.Bytes())
	return nil
}
func (tcpConn *TCPConn) InitCrypt(k []byte) {
	tcpConn.Crypt.Init(k)
}
