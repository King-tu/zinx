package iface

import (
	"net"
)

type IConnect interface {
	Start()
	Stop()
	Send([]byte, uint32) (int, error)
	GetTCPConn() *net.TCPConn
	GetConnId() uint32
	SetProperty(string, interface{})
	GetProperty(string) interface{}
}

//type HandleFunc func(IRequest)

