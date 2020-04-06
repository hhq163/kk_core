package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"mangos/common"
	"mangos/core/auth"
	"mangos/core/slog"
	"mangos/core/util"
	"net"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type WebsocketConnSet map[*websocket.Conn]struct{}

type WSConn struct {
	conn      *websocket.Conn
	writeChan *util.SyncQueue
	maxMsgLen uint32
	closeFlag int32
	Crypt     auth.AuthCrypt
}

func newWSConn(conn *websocket.Conn, maxMsgLen uint32) *WSConn {
	wsConn := new(WSConn)
	wsConn.conn = conn
	wsConn.writeChan = util.NewSyncQueue()
	wsConn.maxMsgLen = maxMsgLen

	go func() {
		for {
			bs := wsConn.writeChan.PopAll()
			if bs == nil {
				break
			}
			for _, b := range bs {
				buf, ok := b.([]byte)
				if ok {
					err := conn.WriteMessage(websocket.BinaryMessage, buf)
					if err != nil {
						slog.Info(err)
						goto closeSocket
					}
				}
			}

		}
	closeSocket:
		atomic.StoreInt32(&wsConn.closeFlag, 1)
		conn.Close()
	}()

	return wsConn
}

func (wsConn *WSConn) Close() {
	wsConn.writeChan.Close()
}

func (wsConn *WSConn) IsClosed() bool {
	return atomic.LoadInt32(&wsConn.closeFlag) == 1
}

func (wsConn *WSConn) Write(b []byte) {
	wsConn.writeChan.Push(b)
}

func (wsConn *WSConn) DirectWrite(b []byte) {
	wsConn.conn.WriteMessage(websocket.BinaryMessage, b)
}

func (wsConn *WSConn) LocalAddr() net.Addr {
	return wsConn.conn.LocalAddr()
}

func (wsConn *WSConn) RemoteAddr() net.Addr {
	return wsConn.conn.RemoteAddr()
}

func (wsConn *WSConn) ReadMsg() (*common.WorldPacket, error) {
	_, b, err := wsConn.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	wsConn.Crypt.DecryptRecv(b[:4])
	msgLen := int(binary.LittleEndian.Uint16(b[:2]))
	opCode := binary.LittleEndian.Uint16(b[2:4])
	if msgLen != len(b) {
		return nil, errors.New("收到ws数据长度错误")
	}
	packet := &common.WorldPacket{}
	packet.Initialize(opCode)
	packet.WriteBytes(b[4:])
	return packet, err
}

// args must not be modified by the others goroutines
func (wsConn *WSConn) WriteMsg(packet *common.WorldPacket) error {
	if wsConn.IsClosed() {
		return errors.New("socket is closed")
	}
	// get len
	msgLen := uint16(packet.Len() + 4)
	header := new(bytes.Buffer)
	binary.Write(header, binary.LittleEndian, msgLen)
	binary.Write(header, binary.LittleEndian, packet.GetOpCode())
	wsConn.Crypt.EncryptSend(header.Bytes())
	binary.Write(header, binary.LittleEndian, packet.Bytes())
	wsConn.Write(header.Bytes())
	return nil
}

func (wsConn *WSConn) DirectWriteMsg(packet *common.WorldPacket) error {
	if wsConn.IsClosed() {
		return errors.New("socket is closed")
	}
	// get len
	msgLen := uint16(packet.Len() + 4)
	header := new(bytes.Buffer)
	binary.Write(header, binary.LittleEndian, msgLen)
	binary.Write(header, binary.LittleEndian, packet.GetOpCode())
	wsConn.Crypt.EncryptSend(header.Bytes())
	binary.Write(header, binary.LittleEndian, packet.Bytes())
	wsConn.DirectWrite(header.Bytes())
	return nil
}

func (wsConn *WSConn) InitCrypt(k []byte) {
	wsConn.Crypt.Init(k)
}