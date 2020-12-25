package common

import (
	"bytes"
	"reflect"
)

//Packet 消息包处理类
type UdpPacket struct {
	buffer *bytes.Buffer
	cmdId  uint16
	userId uint32
}

//Initialize 初始化操作数
func (p *UdpPacket) Initialize(cmdId uint16) {
	p.buffer = new(bytes.Buffer)
	p.cmdId = cmdId
}

func (p *UdpPacket) Init(cmdId uint32, userId uint32) {
	p.buffer = new(bytes.Buffer)
	p.cmdId = uint16(cmdId)
	p.userId = userId
}

//GetOpCode 操作数
func (p *UdpPacket) GetCmd() uint16 {
	return p.cmdId
}

func (p *UdpPacket) GetUserId() uint32 {
	return p.userId
}

func (p *UdpPacket) WriteBytes(b []byte) {
	p.buffer.Write(b)
}

func (p *UdpPacket) WriteData(data interface{}) {
	v := reflect.Indirect(reflect.ValueOf(data))
	encode(p.buffer, v)
}

func (p *UdpPacket) ReadData(data interface{}) error {
	v := reflect.ValueOf(data)
	return decode(p.buffer, v)
}

func (p *UdpPacket) Len() int {
	return p.buffer.Len()
}

func (p *UdpPacket) Bytes() []byte {
	return p.buffer.Bytes()
}

func (p *UdpPacket) GetBuffer() *bytes.Buffer {
	return p.buffer
}

func (p *UdpPacket) SetBuffer(b *bytes.Buffer) {
	p.buffer = b
}
