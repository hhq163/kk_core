package impl

import (
	"game_server/game/proto"
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
	Handler func(*CSession, map[string]interface{})
}

type AgentHandler struct {
	Name    string
	Status  uint8
	Handler func(*agent, map[string]interface{})
}

func init() {
	OpCodeTable[proto.MSG_GET_AMOUNT] = OpCodeHandler{"MSG_GET_AMOUNT", STATUS_AUTHED, (*CSession).HandleGetAmount}

	AgentCodeTable[proto.MSG_HEARTBEAT] = AgentHandler{"MSG_HEARTBEAT", STATUS_NOT_AUTHED, (*agent).HandleHEARTBEAT}
	AgentCodeTable[proto.MSG_CREATER_ROLE] = AgentHandler{"MSG_CREATER_ROLE", STATUS_NOT_AUTHED, (*agent).HandleCreateRole}
	AgentCodeTable[proto.MSG_LOGIN] = AgentHandler{"MSG_LOGIN", STATUS_NOT_AUTHED, (*agent).HandleLogin}

}
