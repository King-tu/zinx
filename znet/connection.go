package znet

import (
	"fmt"
	"github.com/King-tu/zinx/iface"
	"github.com/King-tu/zinx/utils"
	"io"
	"log"
	"net"
	"sync"
)

type Connection struct {
	Conn *net.TCPConn
	ConnID uint32
	IsClosed bool
	//HandleAPI iface.HandleFunc
	//Router iface.IRouter
	RouterMap iface.IMsgHandler
	MsgChan chan []byte
	Server iface.IServer
	Property map[string] interface{}
	PropertyLock sync.RWMutex
}


func NewConnection(conn *net.TCPConn, connID uint32, routerMap iface.IMsgHandler, server iface.IServer) *Connection {

	return &Connection{
		Conn:      conn,
		ConnID:    connID,
		IsClosed:  false,
		//HandleAPI: handleFunc,
		//Router: router,
		RouterMap: routerMap,
		MsgChan: make(chan []byte),
		Server: server,
		Property: make(map[string]interface{}),
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWriter()

	c.Server.CallStartHookFunc(c)
}

func (c *Connection) StartReader() {
	fmt.Println("[Starteader...]")
	defer fmt.Println("[StartReader goroutine exit]")
	defer c.Stop()

	for {
		dp := NewDataPack()
		//第一次，读数据头
		headBuf := make([]byte, dp.GetDataHeadLen())
		//func ReadFull(r Reader, buf []byte) (n int, err error)
		_, err := io.ReadFull(c.Conn, headBuf)
		if err != nil {
			log.Println("conn.Read err: ", err)
			return
		}

		iMsg, err := dp.UnPack(headBuf)
		if err != nil {
			log.Println("db.UnPack err: ", err)
			return
		}

		dataLen := iMsg.GetDataLen()
		if dataLen == 0 {
			fmt.Println("数据长度为0，无需读取")
			continue
		}

		//第二次，读数据
		dataBuf := make([]byte, dataLen)
		_, err = io.ReadFull(c.Conn, dataBuf)
		if err != nil {
			log.Println("conn.Read err: ", err)
			return
		}
		fmt.Printf("Server <== Client, len: %d, buf: %s\n", dataLen, dataBuf)

		iMsg.SetData(dataBuf)
		req := NewRequest(c, iMsg)

		if utils.GlobalConfig.WorkerSize > 0 {
			c.RouterMap.SendReqToWorkQuene(req)
		} else {
			go c.RouterMap.DoMsgHandler(req)
		}

	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[StartWriter...]")
	defer fmt.Println("[StartWriter goroutine exit]")

	for retInfo := range c.MsgChan {

		_, err := c.Conn.Write(retInfo)
		if err != nil {
			log.Println("tcpConn.Write err: ", err)
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.IsClosed {
		return
	}
	c.Server.GetConnMgr().RemoveConn(int(c.ConnID))
	//关闭连接之前调用钩子函数
	c.Server.CallStopHookFunc(c)
	c.IsClosed = true
	close(c.MsgChan)

	_ = c.Conn.Close()
}
func (c *Connection) Send(data []byte, msgId uint32) (int, error) {
	msg := NewMessage(msgId, uint32(len(data)), data)
	dp := NewDataPack()
	retInfo, err := dp.Pack(msg)
	if err != nil {
		log.Println(" dp.Pack err: ", err)
		return -1, err
	}

	c.MsgChan <- retInfo

	return -1, nil
}

func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnID
}

func (c *Connection) GetProperty(key string) interface{} {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	return c.Property[key]
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	c.Property[key] = value
}

func (c *Connection) RemoveProperty(key string) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	delete(c.Property , key)
}

