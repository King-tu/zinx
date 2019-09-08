package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)


func TestDataPack_Pack(t *testing.T) {
	fmt.Println("TestDataPack_Pack called...")
	//server
	go func() {
			listener, err := net.Listen("tcp", ":8888")
			if err != nil {
				t.Error("net.Listen err: ", err)
				return
			}
			defer listener.Close()

			conn, err := listener.Accept()
			if err != nil {
				t.Error("net.Listen err: ", err)
				return
			}
			defer conn.Close()

		for {
			//var buf bytes.Buffer
			headBuf := make([]byte, 8)

			//_, err := conn.Read(&headBuf)
			//func ReadFull(r Reader, buf []byte) (n int, err error)
			_, err := io.ReadFull(conn, headBuf)
			if err != nil {
				t.Error("conn.Read err: ", err)
				return
			}

			db := NewDataPack()
			iMsg, err := db.UnPack(headBuf)
			if err != nil {
				t.Error("db.UnPack err: ", err)
				return
			}

			dataLen := iMsg.GetDataLen()
			if dataLen == 0 {
				fmt.Println("数据长度为0，无需读取")
				continue
			}
			iMsgId := iMsg.GetMsgID()

			dataBuf := make([]byte, dataLen)
			_, err = io.ReadFull(conn, dataBuf)
			if err != nil {
				t.Error("conn.Read err: ", err)
				return
			}
			//iMsg.SetData(dataBuf)

			fmt.Printf("len: %d, msgId: %d, msg: %s\n", dataLen, iMsgId, dataBuf)
		}

	}()

	//client
	go func() {
		data1 := []byte("你好")
		data2 := []byte("hello world")

		msg1 := NewMessage(0, 0, nil)
		msg2 := NewMessage(0, uint32(len(data1)), data1)
		msg3 := NewMessage(0, uint32(len(data2)), data2)


		db := NewDataPack()
		buf, _ := db.Pack(msg1)
		buf2, _ := db.Pack(msg2)
		buf3, _ := db.Pack(msg3)
		buf = append(buf, buf2...)
		buf = append(buf, buf3...)

		conn, err := net.Dial("tcp", "127.0.0.1:8888")
		if err != nil {
			t.Error("net.Dial err: ", err)
			return
		}
		defer conn.Close()

		cnt, err := conn.Write(buf)
		if err != nil {
			t.Error("conn.Write err: ", err)
			return
		}
		fmt.Println("cnt: ", cnt, "len: ", len(buf))
	}()

	time.Sleep(2 * time.Second)

/*	var dp DataPack
	msg := &Message{
		MsgID: 1,
		Len:   10,
		Data:  []byte("helloworld"),
	}

	buf, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("dp.Pack err: ", err)
	}

	fmt.Printf("%d\n", len(buf))


	retMsg, err := dp.UnPack(buf[:8])
	if err != nil {
		fmt.Println("dp.UnPack err: ", err)
	}

	fmt.Println("retMsg len: ", retMsg.GetDataLen())
	fmt.Println("retMsg ID: ", retMsg.GetMsgID())*/
}
