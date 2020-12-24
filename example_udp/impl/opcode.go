package impl

import (
	"udp_server/protocol"

	"github.com/hhq163/kk_core/common"
)

var OpCodeTable = make(map[uint16]OpCodeHandler, 0)

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

func init() {
	// OpCodeTable[uint16(protocol.MSG_GET_AMOUNT)] = OpCodeHandler{"MSG_GET_AMOUNT", STATUS_AUTHED, (*CSession).HandleGetAmount}

	OpCodeTable[uint16(protocol.Cmd_CBeat)] = OpCodeHandler{"Cmd_CBeat", STATUS_NOT_AUTHED, (*CSession).HandleHEARTBEAT}
	OpCodeTable[uint16(protocol.Cmd_CLogin)] = OpCodeHandler{"Cmd_CLogin", STATUS_NOT_AUTHED, (*CSession).HandleLogin}

}
