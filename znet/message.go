package znet

import 	"github.com/King-tu/zinx/iface"

type Message struct {
	MsgID uint32
	Len uint32
	Data []byte
}

func NewMessage(MsgId, dataLen uint32, data []byte) iface.IMessage {
	return &Message{
		MsgID: MsgId,
		Len:   dataLen,
		Data:  data,
	}
}

func (msg *Message) GetMsgID() uint32 {
	return msg.MsgID
}
func (msg *Message) GetDataLen() uint32 {
	return msg.Len
}
func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) SetMsgID(msgId uint32) {
	msg.MsgID = msgId
}
func (msg *Message) SetDataLen(dataLen uint32) {
	msg.Len = dataLen
}
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}