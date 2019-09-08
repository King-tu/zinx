package iface

type IMsgHandler interface {
	AddRouter(uint32, IRouter)
	DoMsgHandler(IRequest)
	StartWorkPool ()
	SendReqToWorkQuene(IRequest)
}
