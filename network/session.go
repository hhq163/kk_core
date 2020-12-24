package network

import "github.com/hhq163/kk_core/common"

type Sess interface {
	Update() bool
	QueuePacket(msg common.IPacket)
	SendPacket(msg common.IPacket)
}
