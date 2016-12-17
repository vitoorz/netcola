package support

import (
	pb "github.com/golang/protobuf/proto"
	"library/core/datamsg"
	"net"
	"service"
	"sync"
	. "types"
)

var MyServerId = IdString("0x84a5d1600002c001")

type AsyncSender struct {
	sync.Mutex
	Objects    []*NetMsg
	msgPool    *service.BufferPool
	serverCoon net.Conn
	serverName string
}

func (s *AsyncSender) SetSyncBuffPool(pool *service.BufferPool) {
	s.msgPool = pool
}

func (s *AsyncSender) SetGsConnection(name string, con net.Conn) {
	s.serverCoon = con
	s.serverName = name
}

func (s *AsyncSender) SendClientNotify(code MsgType, userId IdString, notify pb.Message) {
	s.send(code, NetMsgIdFlagClient, notify, userId)
}

func (s *AsyncSender) SendServerNotify(code MsgType, notify pb.Message) {
	s.send(code, NetMsgIdFlagServer, notify, MyServerId)
}

func (s *AsyncSender) SendBroadCastNotify(code MsgType, groupId IdString, notify pb.Message) {
	s.send(code, NetMsgIdFlagBroadCast, notify, groupId)
}

func (s *AsyncSender) send(code MsgType, flag uint32, payload pb.Message, receiver IdString) {
	sp := &NetMsg{}
	sp.ObjectID = receiver.ToObjectID()
	sp.SetPayLoad(code, payload, flag)

	s.Lock()
	s.Objects = append(s.Objects, sp)
	s.Unlock()
}

func (s *AsyncSender) getNetMessages() []*NetMsg {
	s.Lock()
	objects := s.Objects
	s.Objects = nil
	s.Unlock()

	return objects
}

func (s *AsyncSender) GetAsyncNetMessages() []*NetMsg {
	objects := s.getNetMessages()
	return objects
}

//self motivated
func (s *AsyncSender) SendInstant(code MsgType, flag uint32, receiver IdString, payload pb.Message) {
	sp := &NetMsg{}
	sp.ObjectID = receiver.ToObjectID()
	sp.SetPayLoad(code, payload, flag)

	msg := datamsg.NewDataMsg("", "", DataMsgFlagS2G, nil)
	msg.Payload = sp
	msg.SetMeta(s.serverName, s.serverCoon)

	s.msgPool.Append(msg)
}
