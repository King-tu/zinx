package znet

import 	"github.com/King-tu/zinx/iface"

type Request struct {
	IConn iface.IConnect
	Msg iface.IMessage
}

func NewRequest(iconn iface.IConnect, msg iface.IMessage) iface.IRequest {
	return &Request{
		IConn: iconn,
		Msg: msg,
	}
}

func (req *Request) GetIConn() iface.IConnect {
	return req.IConn
}
func (req *Request) GetMessage() iface.IMessage {
	return req.Msg
}


