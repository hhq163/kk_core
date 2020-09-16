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

		if packet.GetCmd() != protocol.Cmd_CBeat {
			tLog.Debug("packet.GetCmd()=", packet.GetCmd())
		}

		agentHandle := AgentCodeTable[packet.GetCmd()] //未登录业务处理
		if agentHandle.Handler != nil {
			agentHandle.Handler(a, packet.GetCmd())
		} else {

			if a.session != nil && a.auth {
				a.session.QueuePacket(msg)
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
