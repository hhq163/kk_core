package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"sync/atomic"
	"time"

	"github.com/hhq163/kk_core/auth"
	"github.com/hhq163/kk_core/base"
	"github.com/hhq163/kk_core/common"
	"github.com/hhq163/kk_core/util"

	"github.com/gorilla/websocket"
)

const wsmLen uint32 = 4 //长度占用的字节数

type WSConn struct {
	conn       *websocket.Conn
	writeQueue *util.SyncQueue
	maxMsgLen  uint32
	closeFlag  int32
	activeTime int64 //the time of last receive msg
	Crypt      auth.AuthCrypt
}

func NewWSConn(conn *websocket.Conn, maxMsgLen uint32) *WSConn {
	w := new(WSConn)
	w.conn = conn
	w.writeQueue = util.NewSyncQueue()
	w.maxMsgLen = maxMsgLen

	go func() {
		for {
			bs := w.writeQueue.PopAll()
			if bs == nil {
				break
			}
			for _, b := range bs {
				buf, ok := b.([]byte)
				if ok {
					err := conn.WriteMessage(websocket.BinaryMessage, buf)
					if err != nil {
						base.Log.Info(err)
						goto closeSocket
					}
				}
			}

		}
	closeSocket:
		atomic.StoreInt32(&w.closeFlag, 1)
		conn.Close()
	}()

	return w
}

func (w *WSConn) Close() {
	w.writeQueue.Close()
}

func (w *WSConn) IsClosed() bool {
	return atomic.LoadInt32(&w.closeFlag) == 1
}

func (w *WSConn) Write(b []byte) {
	w.writeQueue.Push(b)
}

//不进队列，直接发送
func (w *WSConn) DirectWrite(b []byte) {
	w.conn.WriteMessage(websocket.BinaryMessage, b)
}

func (w *WSConn) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *WSConn) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}

func (w *WSConn) Read(b []byte) (n int, err error) {
	var l int
	_, msg, err := w.conn.ReadMessage()
	if err != nil {
		return l, err
	}
	b = msg
	l = len(msg)
	return l, nil
}

func (w *WSConn) ReadMsg() (common.IPacket, error) {
	_, b, err := w.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	w.Crypt.DecryptRecv(b[:wsmLen+2])
	msgLen := int(binary.LittleEndian.Uint16(b[:wsmLen]))
	cmdId := binary.LittleEndian.Uint16(b[wsmLen:2])
	if msgLen != len(b) {
		return nil, errors.New("收到ws数据长度错误")
	}
	packet := &common.Packet{}
	packet.Initialize(cmdId)
	packet.WriteBytes(b[wsmLen+2:])
	return packet, err
}

func (w *WSConn) WriteMsg(packet common.IPacket) error {
	if w.IsClosed() {
		return errors.New("socket is closed")
	}
	msgLen := uint16(packet.Len() + int(wsmLen) + 2)
	header := new(bytes.Buffer)
	binary.Write(header, binary.LittleEndian, msgLen)
	binary.Write(header, binary.LittleEndian, packet.GetCmd())
	w.Crypt.EncryptSend(header.Bytes())
	binary.Write(header, binary.LittleEndian, packet.Bytes())
	w.Write(header.Bytes())
	return nil
}

func (w *WSConn) InitCrypt(k []byte) {
	w.Crypt.Init(k)
}

func (w *WSConn) IsTimeout(maxTime uint32) bool {
	var ret bool
	cur := uint32(time.Now().Unix() - w.activeTime)
	if cur > maxTime {
		ret = true
	}
	return ret
}
