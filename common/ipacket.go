package common

import "bytes"

//消息接口
type IPacket interface {
	Initialize(cmdId uint16)
	GetCmd() uint16
	WriteBytes(p []byte)
	WriteData(data interface{})
	Len() int
	Bytes() []byte
	GetBuffer() *bytes.Buffer
	SetBuffer(b *bytes.Buffer)
}
