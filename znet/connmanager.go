package znet

import (
	"fmt"
	"sync"
	"github.com/King-tu/zinx/iface"
)

type ConnManager struct {
	Conns map[int]iface.IConnect
	ConnLock sync.RWMutex
}

func NewConnManager() iface.IConnManager {
	return &ConnManager{
		Conns:make(map[int]iface.IConnect),
	}
}

func (cm *ConnManager) AddConn(conn iface.IConnect) {
	fmt.Println("添加链接", conn.GetConnId())
	cid := int(conn.GetConnId())

	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	if _, ok := cm.Conns[cid]; ok {
		fmt.Println("链接已存在，无需添加", cid)
		return
	}

	cm.Conns[cid] = conn
}
func (cm *ConnManager) RemoveConn(cid int) {
	fmt.Println("删除链接", cid)
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	delete(cm.Conns, cid)

}
func (cm *ConnManager) GetConn(cid int) iface.IConnect {
	fmt.Println("获取一个链接")
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	return cm.Conns[cid]
}

func (cm *ConnManager) GetConnCount() int {
	fmt.Println("获取所有链接数量")

	return len(cm.Conns)

}

func (cm *ConnManager) ClearConn() {
	fmt.Println("清除所有链接")

	for cid, conn := range cm.Conns {
		conn.Stop()

		delete(cm.Conns, cid)
	}
}
