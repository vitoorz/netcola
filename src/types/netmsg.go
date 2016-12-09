package types

import (
	"bytes"
	"encoding/binary"
	pb "github.com/golang/protobuf/proto"
	"library/logger"
)

const (
	NetMsgHeadNoIdSize  = 4 + 4
	NetMsgHeadSize      = 8 + 4 + 4
)

type NetMsgHeadNoId struct {
	OpCode   MsgType
	Size     int32
}


type NetMsgHead struct{
	UserId   PlayerId
	NetMsgHeadNoId
}

func (h *NetMsgHead) DecodeHead(head []byte) error {
	r := bytes.NewReader(head)

	//network flow always use BigEndian
	return binary.Read(r, binary.BigEndian, h)
}

func (h *NetMsgHead) EncodeHead() ([]byte, error) {
	w := new(bytes.Buffer)

	err := binary.Write(w, binary.BigEndian, h)

	return w.Bytes(), err
}

func (h *NetMsgHead) DecodeHeadNoId(head []byte) error {
	r := bytes.NewReader(head)

	//network flow always use BigEndian
	return binary.Read(r, binary.BigEndian, &(h.NetMsgHeadNoId))
}

func (h *NetMsgHead) EncodeHeadNoID() ([]byte, error) {
	w := new(bytes.Buffer)

	err := binary.Write(w, binary.BigEndian, &(h.NetMsgHeadNoId))

	return w.Bytes(), err
}

type NetMsg struct {
	NetMsgHead
	Content    []byte
	payload    pb.Message
}

func NewNetMsg(userId PlayerId, code MsgType) *NetMsg {
	msg := &NetMsg{NetMsgHead:NetMsgHead{UserId: userId}}
	msg.OpCode = code;

	return msg
}

func NewNetMsgFromHead(head []byte) (*NetMsg, error){
	msg := &NetMsg{}
	if err := msg.DecodeHead(head); err != nil {
		return nil, err
	}

	msg.Content = make([]byte, msg.Size)
	return msg, nil
}

func NewNetMsgFromHeadNoId(head []byte) (*NetMsg, error){
	msg := &NetMsg{}
	if err := msg.DecodeHeadNoId(head); err != nil {
		return nil, err
	}

	msg.Content = make([]byte, msg.Size)
	return msg, nil
}

func (msg *NetMsg) SetPayLoad(payLoad pb.Message) error {
	content, err := pb.Marshal(payLoad)
	if err != nil {
		logger.Error("marsh type %d payload %+v error : %s", msg.OpCode, payLoad, err.Error())
		return err
	}

	msg.payload = payLoad
	msg.Content = content
	msg.Size    = int32(len(content))

	return nil
}

func (msg *NetMsg) BinaryProto() ([]byte, error) {
	binary, err := msg.EncodeHead()
	if err != nil {
		return nil, err
	}

	binary = append(binary, msg.Content...)

	return binary, nil
}

func (msg *NetMsg) BinaryProtoNoId() ([]byte, error) {
	binary, err := msg.EncodeHeadNoID()
	if err != nil {
		return nil, err
	}

	binary = append(binary, msg.Content...)

	return binary, nil
}