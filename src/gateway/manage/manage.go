package manage

import (
	"errors"
	pb "github.com/golang/protobuf/proto"
	"library/logger"
	"net"
	"sync"
	. "types"
)

const (
	connMetaTypeServer = 1
	connMetaTypeClient = 2
)

var Servers = &metaManage{
	connections: make(map[IdString]*ConnMeta, 0),
	metaType:    connMetaTypeServer,
}

var Clients = &metaManage{
	connections: make(map[IdString]*ConnMeta, 0),
	metaType:    connMetaTypeClient,
}

func NewServerMeta(conn *net.TCPConn) *ConnMeta {
	return &ConnMeta{metaType: connMetaTypeServer, Conn: conn}
}

func NewClientMeta(conn *net.TCPConn) *ConnMeta {
	return &ConnMeta{metaType: connMetaTypeClient, Conn: conn}
}

type ConnMeta struct {
	ID          IdString
	Conn        *net.TCPConn
	LoginAt     UnixTS
	LastSend    UnixTS
	LastRecv    UnixTS
	SentByte    int
	RecvByte    int
	SentNum     int
	RecvNum     int
	ForwardMeta *ConnMeta

	ObjID ObjectID

	metaType int
}

type metaManage struct {
	lock        sync.Mutex
	connections map[IdString]*ConnMeta
	metaType    int
}

func (mm *metaManage) mineMeta(meta *ConnMeta) bool {
	if meta == nil {
		return false
	}
	if meta.metaType != mm.metaType {
		logger.Error("Connection meta useage error: mm %d accept %d", mm.metaType, meta.metaType)
		return false
	}
	return true
}

func (mm *metaManage) GetMeta(id IdString) (*ConnMeta, bool) {
	mm.lock.Lock()
	meta, ok := mm.connections[id]
	mm.lock.Unlock()
	return meta, ok
}

func (mm *metaManage) Logout(meta *ConnMeta) {
	if ok := mm.mineMeta(meta); !ok {
		return
	}

	mm.lock.Lock()
	delete(mm.connections, meta.ID)
	mm.lock.Unlock()

	if meta.Conn != nil {
		meta.Conn.Close()
		meta.Conn = nil
	}
	meta.ForwardMeta = nil
}

func (mm *metaManage) Login(meta *ConnMeta, content []byte) bool {
	if ok := mm.mineMeta(meta); !ok {
		return false
	}

	switch meta.metaType {
	case connMetaTypeClient:
		req := &LoginReq{}
		if err := pb.Unmarshal(content, req); err != nil {
			return false
		}
		meta.ID = IdString(req.UserId)
		serverMeta, ok := Servers.GetMeta(IdString(req.ServerId))
		if !ok || serverMeta.Conn == nil {
			logger.Error("client %s login server %s error, server not online", meta.ID, req.ServerId)
			return false
		}
		meta.ForwardMeta = serverMeta

	case connMetaTypeServer:
		req := &ServerLoginReq{}
		if err := pb.Unmarshal(content, req); err != nil {
			return false
		}
		meta.ID = IdString(req.ServerId)
	}

	meta.ObjID = meta.ID.ToObjectID()

	mm.lock.Lock()
	mm.connections[meta.ID] = meta
	mm.lock.Unlock()
	return true
}

func (meta *ConnMeta) Send(b []byte) (int, error) {
	if meta.Conn == nil {
		return 0, errors.New("ConnMeta.Conn is nil")
	}

	n, err := meta.Conn.Write(b)
	meta.SentByte += n
	meta.SentNum += 1

	//todo
	return n, err
}
