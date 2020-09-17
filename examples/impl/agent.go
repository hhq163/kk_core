package impl

import (
	"bytes"
	"encoding/binary"
	"kk_server/base"
	"kk_server/protocol"
	"net"

	"github.com/hhq163/kk_core/common"
	"github.com/hhq163/kk_core/network"
	"github.com/hhq163/logger"
	"github.com/golang/protobuf/proto"
)

type agent struct {
	conn    network.Conn
	session *CSession
	auth    bool
}

//SendPacket send msg
func SendPacket(conn network.Conn, msg common.IPacket) error {

	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, msg.Len())
	binary.Write(data, binary.LittleEndian, msg.GetCmd())
	binary.Write(data, binary.LittleEndian, msg.Bytes())
	conn.Write(data.Bytes())

	// base.Log.Debug("SendPacket msg.MsgType=", msg.MsgType)
	conn.Write(data.Bytes())
	return nil
}

//Run 在socket的recv gorouting执行
func (a *agent) Run() {
	tLog := base.Log.With(logger.FuncName, "agent.Run()")
	tLog.Debug("start")

	defer func() {
		if p := recover(); p != nil {
			tLog.Error("socket panic err: ", p)
		}
	}()

	for {
		packet, err := a.conn.ReadMsg()
		if err != nil {
			tLog.Error("err=", err.Error())
			a.Close()
			break
		}

		if packet.GetCmd() != uint16(protocol.Cmd_CBeat) {
			tLog.Debug("packet.GetCmd()=", packet.GetCmd())
		}

		agentHandle := AgentCodeTable[packet.GetCmd()] //未登录业务处理
		if agentHandle.Handler != nil {
			agentHandle.Handler(a, packet)
		} else {

			if a.session != nil && a.auth {
				a.session.QueuePacket(packet)
			} else {
				tLog.Error("session is nil")
				a.conn.Close()
				return
			}
		}

	}
}

func (a *agent) OnClose() {
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

//用户心跳
func (a *agent) HandleHEARTBEAT(requestMsg common.IPacket) {
	tLog := base.Log.With(logger.FuncName, "Handle_HEARTBEAT()")
	tLog.Debug("in")

	rsp := &common.Packet{}
	rsp.Initialize(uint16(protocol.Cmd_SBeat))
	SendPacket(a.conn, rsp)

	tLog.Debug("out")
}

func (a *agent) HandleLogin(reqMsg common.IPacket) {
	tLog := base.Log.With(logger.FuncName, "Handle_HEARTBEAT()")
	tLog.Debug("in")

	pbData := &protocol.ClientLogin{}
	err := proto.Unmarshal(reqMsg.Bytes(), pbData)
	if err != nil {
		tLog.Error("proto.Unmarshal failed cmd=", reqMsg.GetCmd())
		return // 跳出循环，进行下一次消息读取
	}

	tLog.Debug("out")
}
