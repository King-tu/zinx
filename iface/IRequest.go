package iface

type IRequest interface {
	GetIConn() IConnect
	GetMessage() IMessage
}
