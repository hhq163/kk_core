package network

import (
	"mangos/common"
	"net"
)

type Conn interface {
	ReadMsg() (*common.WorldPacket, error)
	WriteMsg(packet *common.WorldPacket) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	IsClosed() bool
	InitCrypt(k []byte)
	DirectWriteMsg(packet *common.WorldPacket) error
}
