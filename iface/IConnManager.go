package iface

type IConnManager interface {
	AddConn(IConnect)
	RemoveConn(int)
	GetConn(int) IConnect
	GetConnCount() int
	ClearConn()
}
