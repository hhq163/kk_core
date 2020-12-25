package impl

import (
	"sync"
	"time"
	"udp_server/base"

	"github.com/hhq163/kk_core/common"
	"github.com/hhq163/kk_core/network"
	"github.com/hhq163/kk_core/util"
	"github.com/hhq163/logger"
)

const MaxTime = 2 * 60 //保活最大间隔2分钟

type CSession struct {
	UserId     int64
	recvQueue  *util.SyncQueue
	conn       *network.UDPVconn
	logoutTime int64
	once       sync.Once
}

func NewCSession(udpConn *network.UDPVconn) *CSession {
	cs := &CSession{
		recvQueue: util.NewSyncQueue(),
		conn:      udpConn,
	}
	return cs
}

func (s *CSession) SendPacket(msg common.IPacket) {
	s.conn.WriteMsg(msg)
}

func (s *CSession) Update() bool {
	pcks, ok := s.recvQueue.TryPopAll()
	if !ok || pcks == nil {
		goto check
	}
	for _, pck := range pcks {
		msg, ok := pck.(*common.Packet)
		if !ok {
			base.Log.Info("pck.(*proto.C2SMessage) not ok")
			break
		}
		s.handler(msg)
	}
check:
	if s.conn.IsClosed() && s.logoutTime == 0 { //socket is closed
		base.Log.Debug("should be closed()111")
		s.logoutTime = time.Now().Unix()
	}

	///- If necessary, logout
	currTime := time.Now().Unix()
	if s.ShouldLogOut(currTime) {
		base.Log.Debug("should be closed()222")
		return false // Will remove this session from session map
	}

	return true
}

func (s *CSession) ShouldLogOut(curTime int64) bool {
	var ret bool
	if s.logoutTime > 0 && curTime >= s.logoutTime {
		base.Log.Debug("ShouldLogOut() curTime=", curTime, "s.logoutTime=", s.logoutTime)
		ret = true
	}
	if s.conn.IsTimeout(MaxTime) {
		s.Close()
		base.Log.Debug("ShouldLogOut() s.conn.IsTimeout")
		ret = true
	}
	return ret
}

//延迟关闭session
func (s *CSession) CloseSession() {
	//防止出现前面的消息还没回复就断链了.
	timer := time.NewTimer(time.Second * 2)
	for {
		select {
		case <-timer.C:
			s.Close()
			return
		}
	}
}

//连接关闭
func (s *CSession) Close() {
	s.once.Do(func() {
		s.conn.Close()
	})
}

func (s *CSession) handler(msg common.IPacket) {
	tLog := base.Log.With(logger.FuncName, "(s *CSession) handler()")
	defer func() {
		if p := recover(); p != nil {
			tLog.Error("CSession panic err: ", p)
		}
	}()
	tLog.Debug("CSession handler() msg.GetCmd()=", msg.GetCmd(), ",userid=", s.UserId)
	opHandle := OpCodeTable[msg.GetCmd()]
	if opHandle.Handler != nil {
		opHandle.Handler(s, msg)
	} else {
		tLog.Error("unknown opcode", msg.GetCmd())
	}
}

//QueuePacket 消息入队
func (s *CSession) QueuePacket(msg common.IPacket) {
	s.recvQueue.Push(msg)
}

func (s *CSession) HandleGetAmount(packet common.IPacket) {

}

func (s *CSession) HandleHEARTBEAT(packet common.IPacket) {

}

func (s *CSession) HandleLogin(packet common.IPacket) {

}
