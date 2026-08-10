package matching

import (
	"time"

	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/models"
	matchingrepo "github.com/xyy0411/blog/repositories/matching"
)

var (
	MatchHub *Hub
)

type Hub struct {
	clients map[int64]*Client // 连接的客户端

	broadcast chan []byte // 广播消息的通道

	register chan *userClient // 注册请求的通道

	unregister chan int64 // 注销请求的通道

	match chan models.Matching // 匹配请求的通道

	quit chan int64 // 用户主动退出匹配的通道
}

type userClient struct {
	id       int64
	client   *Client
	userName string // 用户名
}

func NewMatchingHub() *Hub {
	MatchHub = &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *userClient),
		unregister: make(chan int64),
		clients:    make(map[int64]*Client),
		match:      make(chan models.Matching),
	}
	return MatchHub
}

func (h *Hub) Run() {
	global.Logger.Info("匹配系统启动")
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client.client
			event := models.MatchEvent{
				Type:    "queueing",
				SelfID:  client.id,
				Message: "匹配中",
				Code:    200,
			}
			sendEvent(client.client, event)
		case info := <-h.match:
			matchedList.MatchUsers(info)
		case id := <-h.unregister:
			client, ok := h.clients[id]
			if ok {
				// 检查用户是否还在匹配队列中（异常断连的情况）
				if _, inQueue := matchedList.matchedList.Load(id); inQueue {
					matchedList.RemoveUserFromQueue(id)

					// 获取用户信息用于保存记录
					repo := matchingrepo.NewRepo(global.DB)
					user, err := repo.GetByUserID(id)
					if err != nil {
						global.Logger.Errorf("获取用户信息失败: %v", err)
						user = models.Matching{UserID: id, UserName: ""}
					}

					// 计算匹配持续时间（秒）
					duration := int(time.Now().Sub(client.connectedAt).Seconds())

					// 保存匹配申请记录（网络错误/异常断连）
					matchedList.saveMatchingApplication(id, user.UserName, false, duration, "", models.ExitReasonError)
				}
				delete(h.clients, id)
			}
		case id := <-h.quit:
			client, _ := h.clients[id]
			event := models.MatchEvent{
				Type:    "cancelled",
				SelfID:  id,
				PeerID:  0,
				Message: "已成功退出匹配",
				Code:    200,
			}
			sendEvent(client, event)
			delete(h.clients, id)
			close(client.send)
		}
	}
}
