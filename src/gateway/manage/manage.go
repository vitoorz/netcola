package manage

import (
	pb "github.com/golang/protobuf/proto"
	"library/logger"
	"net"
	"sync"
	. "types"
)

type ServerConnectionMeta struct {
	ServerId   ServerId
	ServerConn *net.TCPConn
	//ClientsConn map[PlayerId]*net.TCPConn
}

type serverManage struct {
	Servers map[ServerId]*ServerConnectionMeta
	sync.Mutex
}

var servers = &serverManage{
	Servers: make(map[ServerId]*ServerConnectionMeta, 0),
}

func serverLogin(connMeta *ServerConnectionMeta, heartBeat *NetMsg) {
	req := &HeadBeat{}
	if err := pb.Unmarshal(heartBeat.Content, req); err != nil {
		logger.Warn("unmarshal server login req error")
		return
	}
	connMeta.ServerId = ServerId(heartBeat.UserId)
	logger.Info("Server %s logined", connMeta.ServerId.ToIdString())

	//check req.Code == 0
	servers.Lock()
	servers.Servers[connMeta.ServerId] = connMeta
	servers.Unlock()

	binary, _ := heartBeat.BinaryProto()
	connMeta.ServerConn.Write(binary)
}

func ServerLogout(connMeta *ServerConnectionMeta) {
	servers.Lock()
	delete(servers.Servers, connMeta.ServerId)
	servers.Unlock()

	if connMeta.ServerConn != nil {
		connMeta.ServerConn.Close()
	}

	connMeta.ServerConn = nil
}

func getServerMeta(serverId ServerId) (*ServerConnectionMeta, bool) {
	m, ok := servers.Servers[serverId]
	return m, ok
}

type ClientConnectionMeta struct {
	UserId     PlayerId
	UUID       string
	MessageNum int
	LoginAt    UnixTS
	LastBeat   UnixTS

	UserIdBinary []byte
	ClientConn   *net.TCPConn

	ServerMeta *ServerConnectionMeta
}

type clientManage struct {
	Clients map[PlayerId]*ClientConnectionMeta
	sync.Mutex
}

var clients = &clientManage{
	Clients: make(map[PlayerId]*ClientConnectionMeta, 0),
}

func clientLogin(conMeta *ClientConnectionMeta, loginReq *NetMsg) bool {
	req := &LoginReq{}
	if err := pb.Unmarshal(loginReq.Content, req); err != nil {
		logger.Error("unmarshal client login req error")
		return false
	}

	conMeta.UserId = IdString(req.UserId).ToPlayerId()

	serverMeta, ok := getServerMeta(IdString(req.ServerId).ToServerId())
	if !ok {
		logger.Warn("player %s try to login server %s, server not exist", req.UserId, req.ServerId)
		ClientLogout(conMeta)
		return false
	}

	logger.Warn("player %s try to login server %s", req.UserId, req.ServerId)
	conMeta.ServerMeta = serverMeta

	clients.Lock()
	clients.Clients[conMeta.UserId] = conMeta
	clients.Unlock()

	return true
}

func ClientLogout(conMeta *ClientConnectionMeta) bool {

	clients.Lock()
	delete(clients.Clients, conMeta.UserId)
	clients.Unlock()

	if conMeta.ClientConn != nil {
		conMeta.ClientConn.Close()
	}

	conMeta.ClientConn = nil

	return true
}

func getClientMeta(userId PlayerId) (*ClientConnectionMeta, bool) {
	m, ok := clients.Clients[userId]
	if ok {
		return m, true
	}
	return nil, false
}
