package service

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"time"
)

const (
	ScCmdCommonUserStub = 0

	ServiceStatusNotBegin = 0
	ServiceStatusRunning  = 1
	ServiceStatusStopping = 2
	ServiceStatusStopped  = 3
)

type ServiceCmdHandler func(service interface{}, code int, payLoad interface{}) (int, interface{})

type ServiceCommon struct {
	tickStepNano    int64
	handlers        map[int]ServiceCmdHandler
	sysHandlers     map[int]ServiceCmdHandler
	pipe            *dm.DataMsgPipe
	status          int
	serviceInstance interface{} //real service interface data struct pointer
	cm.ControlMsgPipe
}

func NewServiceCommon(tickStepNano int64, pipe *dm.DataMsgPipe) *ServiceCommon {
	if tickStepNano < 1e7 { //10ms at lease
		tickStepNano = 1e7
	}
	if pipe == nil {
		pipe = dm.NewDataMsgPipe(1, 1)
	}
	sc := &ServiceCommon{
		tickStepNano:    tickStepNano,
		handlers:        make(map[int]ServiceCmdHandler),
		sysHandlers:     make(map[int]ServiceCmdHandler),
		pipe:            pipe,
		status:          ServiceStatusNotBegin,
		serviceInstance: nil,
	}
	sc.ControlMsgPipe = *cm.NewControlMsgPipe()
	return sc
}

func (sc *ServiceCommon) SetRealInstance(realService interface{}) {
	sc.serviceInstance = realService
}

func (sc *ServiceCommon) Status() int {
	return sc.status
}

func (sc *ServiceCommon) SysHandler(code int, handler ServiceCmdHandler) {
	sc.sysHandlers[code] = handler
}

func (sc *ServiceCommon) HandleSysCmd(cmd *cm.ControlMsg) {
	sc.Cmd <- cmd
}

func (sc *ServiceCommon) UseHandler(code int, handler ServiceCmdHandler) {
	sc.handlers[code] = handler
}

func (sc *ServiceCommon) Handle(msg *dm.DataMsg) {
	sc.pipe.WriteDownChanNB(msg)
}

func (sc *ServiceCommon) Start(pipe *dm.DataMsgPipe) {
	if pipe != nil {
		sc.pipe = pipe
	}

	go func() {
		sc.status = ServiceStatusRunning
		tickChan := time.NewTicker(time.Duration(sc.tickStepNano)).C
		for sc.status == ServiceStatusRunning {
			select {
			case <-tickChan:
				sc.handleSysCmd(cm.ControlMsgTick, time.Now().UnixNano())
			case cmd, ok := <-sc.Cmd:
				if !ok {
					continue
				}
				sc.handleSysCmd(cmd.MsgType, cmd.Payload)
			case userMsg, okUser := <-sc.pipe.ReadDownChan():
				if !okUser {
					continue
				}
				sc.handleUserCmd(userMsg.MsgType, userMsg)
			default:
			}
		}
		sc.status = ServiceStatusStopped
	}()
}

func (sc *ServiceCommon) Stop() {
	sc.Cmd <- &cm.ControlMsg{cm.ControlMsgExit, nil}
}

func (sc *ServiceCommon) handleSysCmd(code int, payLoad interface{}) (int, interface{}) {
	if code == cm.ControlMsgExit {
		sc.status = ServiceStatusStopping
	}

	handler, ok := sc.sysHandlers[code]
	if !ok {
		return 0, nil
	}
	return handler(sc.serviceInstance, code, payLoad)
}

func (sc *ServiceCommon) handleUserCmd(code int, payLoad interface{}) (int, interface{}) {
	handler, ok := sc.handlers[code]
	if !ok && code >= 0 {
		if handler, ok = sc.handlers[ScCmdCommonUserStub]; !ok {
			return 0, nil
		}
	}
	return handler(sc.serviceInstance, code, payLoad)
}
