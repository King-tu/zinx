package znet

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"github.com/King-tu/zinx/iface"
	"github.com/King-tu/zinx/utils"
)

type Server struct {
	IP        string
	Port      uint32
	Name      string
	IPVersion string
	//Router    iface.IRouter
	RouterMap iface.IMsgHandler
	ConnManager iface.IConnManager

	StartHookFunc func(iface.IConnect)
	StoptHookFunc func(iface.IConnect)
}


func NewServer() iface.IServer {
	return &Server{
		IP:        utils.GlobalConfig.IP,
		Port:      utils.GlobalConfig.Port,
		Name:      utils.GlobalConfig.Name,
		IPVersion: utils.GlobalConfig.IPVerson,
		//Router:    &Router{},
		RouterMap: NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
}

func (s *Server) Start()  {
	fmt.Println("Server start...")
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	address := fmt.Sprintf("%s:%d", s.IP, s.Port)
	tcpAddress, err := net.ResolveTCPAddr(s.IPVersion, address)
	if err != nil {
		log.Println("znet.ResolveTCPAddr err: ", err)
		return
	}

	tcpListener, err := net.ListenTCP(s.IPVersion, tcpAddress)
	if err != nil {
		log.Println("znet.ListenTCP err: ", err)
		return
	}
	//启动线程池
	s.RouterMap.StartWorkPool()

	var connID uint32
	connID = 0

	go func() {
		for {
			//fmt.Println("等待客户端连接：")
			tcpConn, err := tcpListener.AcceptTCP()
			if err != nil {
				log.Println("tcpListener.AcceptTCP err: ", err)
				//return
				continue
			}

			if s.ConnManager.GetConnCount() >= utils.GlobalConfig.MaxConnNum {
				fmt.Println("已达到最大链接数，当前链接被拒绝, 链接id：", connID)
				tcpConn.Close()
				continue
			}

			go func() {

				connection := NewConnection(tcpConn, connID, s.RouterMap, s)
				s.ConnManager.AddConn(connection)
				connID++

				go connection.Start()
			}()
		}
	}()
}

func (s *Server) Stop()  {
	fmt.Println("Server stop...")
//	释放资源
	s.ConnManager.ClearConn()
}

func (s *Server) Serve()  {

	s.Start()
	fmt.Println("Server serve...")

	for{
		runtime.GC()
	}
}

func (s *Server) AddRouter(msgId uint32, router iface.IRouter) {
	s.RouterMap.AddRouter(msgId, router)
}

func (s *Server) GetConnMgr() iface.IConnManager {
	return s.ConnManager
}

func (s *Server) RegisterStartHookFunc(f func(iface.IConnect))  {
	s.StartHookFunc = f

}

func (s *Server) RegisterStopHookFunc(f func(iface.IConnect))  {
	s.StoptHookFunc = f
}

func (s *Server) CallStartHookFunc(conn iface.IConnect)  {
	if s.StartHookFunc != nil {
		s.StartHookFunc(conn)
	}
}

func (s *Server) CallStopHookFunc(conn iface.IConnect)  {
	if s.StoptHookFunc != nil {
		s.StoptHookFunc(conn)
	}
}