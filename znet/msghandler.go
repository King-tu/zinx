package znet

import (
	"fmt"
	"log"
	"zinxProject/v11-connproperty/zinx/iface"
	"zinxProject/v11-connproperty/zinx/utils"
)

type MsgHandler struct {
	RouterMap map[uint32] iface.IRouter
	WorkSize int
	TaskQuene []chan iface.IRequest
}

func NewMsgHandler() iface.IMsgHandler {
	workSize := utils.GlobalConfig.WorkerSize
	return &MsgHandler{
		RouterMap:make(map[uint32] iface.IRouter),
		WorkSize: workSize,
		TaskQuene: make([]chan iface.IRequest, workSize),
	}
}

func (mh *MsgHandler) StartWorkPool () {
	fmt.Println("[StarWorker Pool]...")

	for i := 0; i < utils.GlobalConfig.WorkerSize; i++ {
		fmt.Println("启动worker, worker id:", i)

		go func(i int) {

			mh.TaskQuene[i] = make(chan iface.IRequest, utils.GlobalConfig.TaskQueneSize)
			for {
				req := <- mh.TaskQuene[i]
				mh.DoMsgHandler(req)
			}
		}(i)
	}
}

func (mh *MsgHandler) SendReqToWorkQuene(req iface.IRequest)  {
	cid := req.GetIConn().GetConnId()

	workerId := int(cid) % utils.GlobalConfig.WorkerSize

	fmt.Println("添加cid:", cid, " 的请求到workerid:", workerId)

	mh.TaskQuene[workerId] <- req
}

func (mh *MsgHandler) AddRouter(msgId uint32, router iface.IRouter)  {
	if _, ok := mh.RouterMap[msgId]; !ok {
		mh.RouterMap[msgId] = router
	}
}

func (mh *MsgHandler) DoMsgHandler(req iface.IRequest) {
	msgId := req.GetMessage().GetMsgID()
	router, ok := mh.RouterMap[msgId]
	if !ok {
		log.Printf("不存在msgid=%d 对应的路由!", msgId)
		return
	}

	router.PreHandle(req)
	router.Handle(req)
	router.PostHandle(req)
}
