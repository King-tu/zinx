package iface

type IMessage interface {
	GetMsgID() uint32
	GetDataLen() uint32
	GetData() []byte

	SetMsgID(msgId uint32)
	SetDataLen(dataLen uint32)
	SetData(data []byte)
}
