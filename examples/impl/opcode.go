package impl

import (
	"github.com/hhq163/kk_core/common"
)

var OpCodeTable = make(map[string]OpCodeHandler, 0)
var AgentCodeTable = make(map[string]AgentHandler, 0)

const (
	STATUS_NOT_AUTHED = 0
	STATUS_AUTHED     = 1
	STATUS_UNHANDLED  = 2
)

//消息路由器
type OpCodeHandler struct {
	Name    string
	Status  uint8
	Handler func(*CSession, common.IPacket)
}

type AgentHandler struct {
	Name    string
	Status  uint8
	Handler func(*agent, common.IPacket)
}

func init() {
	OpCodeTable[protocol.MSG_GET_AMOUNT] = OpCodeHandler{"MSG_GET_AMOUNT", STATUS_AUTHED, (*CSession).HandleGetAmount}

	AgentCodeTable[proto.Cmd_CBeat] = AgentHandler{"Cmd_CBeat", STATUS_NOT_AUTHED, (*agent).HandleHEARTBEAT}
	AgentCodeTable[protocol.Cmd_CLogin] = AgentHandler{"Cmd_CLogin", STATUS_NOT_AUTHED, (*agent).HandleLogin}

}
