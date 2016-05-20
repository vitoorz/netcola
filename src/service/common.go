package service

import (
	"library/netmsg"
	"time"
)

const (
	ScCmdCommonUserStub    = 0
	ScCmdStopService       = -1
	ScCmdSetUserCmdHandler = -2
	ScCmdSetTickHandler    = -3
	ScCmdOnTick            = -4

	ServiceStatusNotBegin = 0
	ServiceStatusRunning  = 1
	ServiceStatusStoping  = 2
	ServiceStatusStoped   = 3
)

type ServiceCmd struct {
	Code    int
	PayLoad interface{}
}

type ServiceCmdHandler func(service interface{}, code int, payLoad interface{}) (int, interface{})

type ServiceCommon struct {
	tickStepNano    int64
	handlers        map[int]ServiceCmdHandler
	sysCmd          chan ServiceCmd
	pipe            *netmsg.NetMsgPipe
	status          int
	serviceInstance interface{} //real service interface data struct pointer
}

func NewServiceCommon(tickStepNano int64, pipe *netmsg.NetMsgPipe) *ServiceCommon {
	if tickStepNano < 1e7 { //10ms at lease
		tickStepNano = 1e7
	}
	if pipe == nil {
		pipe = netmsg.NewNetMsgPipe(1, 1)
	}
	return &ServiceCommon{
		tickStepNano:    tickStepNano,
		handlers:        make(map[int]ServiceCmdHandler),
		sysCmd:          make(chan ServiceCmd),
		pipe:            pipe,
		status:          ServiceStatusNotBegin,
		serviceInstance: nil,
	}
}

func (sc *ServiceCommon) SetRealInstance(realService interface{}) {
	sc.serviceInstance = realService
}

func (sc *ServiceCommon) Status() int {
	return sc.status
}

func (sc *ServiceCommon) UseHandler(code int, handler ServiceCmdHandler) {
	sc.handlers[code] = handler
}

func (sc *ServiceCommon) Handle(msg *netmsg.NetMsg) {
	sc.pipe.WriteRecvChan(msg)
}

func (sc *ServiceCommon) HandleServiceCmd(cmd ServiceCmd) {
	sc.sysCmd <- cmd
}

func (sc *ServiceCommon) Start(pipe *netmsg.NetMsgPipe) {
	if pipe != nil {
		sc.pipe = pipe
	}

	go func() {
		sc.status = ServiceStatusRunning
		tickChan := time.NewTicker(time.Duration(sc.tickStepNano)).C
		for sc.status == ServiceStatusRunning {
			select {
			case <-tickChan:
				sc.handleCmd(ScCmdOnTick, time.Now().UnixNano())
			case cmd, ok := <-sc.sysCmd:
				if !ok {
					continue
				}
				sc.handleCmd(cmd.Code, cmd.PayLoad)
			case userMsg, okUser := <-*sc.pipe.ReadRecvChan():
				if !okUser {
					continue
				}
				sc.handleCmd(userMsg.OpCode(), userMsg)
			default:
			}
		}
		sc.status = ServiceStatusStoped
	}()
}

func (sc *ServiceCommon) Stop() {
	sc.status = ServiceStatusStoping
}

func (sc *ServiceCommon) handleCmd(code int, payLoad interface{}) (int, interface{}) {
	handler, ok := sc.handlers[code]
	if !ok && code >= 0 {
		if handler, ok = sc.handlers[ScCmdCommonUserStub]; !ok {
			return 0, nil
		}
	}
	return handler(sc.serviceInstance, code, payLoad)
}
