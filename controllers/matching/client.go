package matching

import (
	"github.com/gorilla/websocket"
	"github.com/xyy0411/blog/global"
	"time"
)

type Client struct {
	hub        *Hub
	conn       *websocket.Conn
	send       chan []byte // 发送消息的通道
	limitTimer chan int64  // 限制匹配时间的通道
	close      chan bool   // 关闭通道
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:        hub,
		conn:       conn,
		send:       make(chan []byte, 256),
		limitTimer: make(chan int64),
		close:      make(chan bool),
	}
}

func (c *Client) checkLimitTimer(id int64) {
	timer := time.NewTimer(0)
	timer.Stop()
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

	for {
		select {
		case <-c.close:
			global.Logger.Infof("用户:%d 主动退出匹配队列", id)
			return
		case t := <-c.limitTimer:
			// 创建或重置定时器
			timer = time.NewTimer(time.Duration(t) * time.Second)
		case <-timer.C:
			matchedList.RemoveUserFromQueue(id)
			c.send <- []byte("匹配超时,已退出匹配队列")
			return
		}
	}
}

func (c *Client) writePump(userID int64) {
	defer func() {
		c.hub.unregister <- userID
		c.conn.Close()
		/*global.Logger.Infof("已与用户:%d 断开连接", userID)*/
	}()

	count := 0
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					global.Logger.Error(err)
					return
				}
				return
			}
			count++
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				global.Logger.Error(err)
				return
			}
			if count == 2 {
				return
			}
		}
	}
}
