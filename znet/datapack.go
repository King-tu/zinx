package znet

import (
	"bytes"
	"encoding/binary"
	"github.com/King-tu/zinx/iface"
)

type DataPack struct {

}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetDataHeadLen() uint32 {
	return 8
}
func (dp *DataPack) Pack(iMsg iface.IMessage) ([]byte, error) {

	len := iMsg.GetDataLen()
	id := iMsg.GetMsgID()
	data := iMsg.GetData()
	var buf bytes.Buffer
	//消息长度
	err := binary.Write(&buf, binary.LittleEndian, len)
	if err != nil {
		return nil, err
	}
	//消息类型
	err = binary.Write(&buf, binary.LittleEndian, id)
	if err != nil {
		return nil, err
	}
	//消息内容
	err = binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		return nil, err
	}

	//fmt.Println("Pack: ", len, id, data)
	return buf.Bytes(), nil
}

//在connection中使用时，会读取两次
//1. 第一次会读取固定8字节的长度，然后调用Unpack
//2. 第二次会读取真实长度的数据
func (dp *DataPack) UnPack(data []byte) (iface.IMessage, error) {

	var msg Message
	reader := bytes.NewReader(data)
	//读消息长度
	err := binary.Read(reader, binary.LittleEndian, &msg.Len)
	if err != nil {
		return nil, err
	}
	//读消息类型
	err = binary.Read(reader, binary.LittleEndian, &msg.MsgID)
	if err != nil {
		return nil, err
	}

	//fmt.Println("UnPack ==== : ", msg.Len, msg.MsgID)

	return &msg, nil
}