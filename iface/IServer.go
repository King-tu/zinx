package iface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(uint32, IRouter)
	GetConnMgr() IConnManager

	RegisterStartHookFunc(func(IConnect))
	RegisterStopHookFunc(func(IConnect))

	CallStartHookFunc(IConnect)
	CallStopHookFunc(IConnect)
}
