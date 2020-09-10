package impl

import (
	"encoding/json"
	"game_server/game/proto"
	"kk_server/base"
	"net"

	"github.com/hhq163/kk_core/network"
	"github.com/hhq163/logger"
)

type agent struct {
	conn    network.Conn
	session *CSession
	auth    bool
}

//SendPacket send msg
func SendPacket(conn network.Conn, msg *proto.S2CMessage) error {
	packet, err := json.Marshal(msg)
	if err != nil {
		base.Log.Error("SendPacket msg.MsgType=", msg.MsgType)
		return err
	}
	// base.Log.Debug("SendPacket msg.MsgType=", msg.MsgType)
	conn.Write(packet)
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
	var packet = make([]byte, base.Setting.Server.MaxMsgLen)

	for {
		n, err := a.conn.Read(packet)
		if err != nil {
			tLog.Error("err=", err.Error())
			a.Close()
			break
		}

		msg := &proto.C2SMessage{}
		err = json.Unmarshal(packet[:n], msg)
		if err != nil {
			tLog.Error("json.Unmarshal error, err=", err.Error())
			continue
		}
		if msg.MsgType != proto.MSG_HEARTBEAT {
			tLog.Debug("msg.MsgType=", msg.MsgType)
		}

		agentHandle := AgentCodeTable[msg.MsgType] //未登录业务处理
		if agentHandle.Handler != nil {
			agentHandle.Handler(a, msg.MsgData)
		} else {

			if a.session != nil && a.auth {
				a.session.QueuePacket(msg)
			} else {
				tLog.Error("session is nil")
				a.conn.Close()
				return
			}
		}

		copy(packet, packet[n:])
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
