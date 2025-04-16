package matching

import (
	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/models"
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
	id     int64
	client *Client
}

func NewMatchingHub() {
	MatchHub = &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *userClient),
		unregister: make(chan int64),
		clients:    make(map[int64]*Client),
		match:      make(chan models.Matching),
	}
	go MatchHub.Run()
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client.client
			client.client.send <- []byte("匹配中")
		case info := <-h.match:
			global.Logger.Info("开始")
			matchedList.MatchUsers(info)
		case id := <-h.unregister:
			if _, ok := h.clients[id]; ok {
				delete(h.clients, id)
			}
		case id := <-h.quit:
			client, _ := h.clients[id]
			client.send <- []byte("已成功退出匹配")
			delete(h.clients, id)
			close(client.send)
		}
	}
}
