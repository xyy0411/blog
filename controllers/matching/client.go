package matching

import (
	"github.com/gorilla/websocket"
	"github.com/xyy0411/blog/global"
	"time"
)

type Client struct {
	hub        *Hub
	conn       *websocket.Conn
	send       chan message // 发送消息的通道
	limitTimer chan int64   // 限制匹配时间的通道
	close      chan bool    // 关闭通道
}

type message struct {
	messageType int8
	message     []byte
}

func newMessage(messageType int8, msg []byte) message {
	return message{
		messageType: messageType,
		message:     msg,
	}
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:        hub,
		conn:       conn,
		send:       make(chan message, 256),
		limitTimer: make(chan int64),
		close:      make(chan bool),
	}
}

func (c *Client) handleTimeout(id int64) {
	timer := time.NewTimer(0)
	timer.Stop()
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

	for {
		select {
		case t := <-c.close:
			if t {
				global.Logger.Infof("用户:%d 意外退出匹配队列,关闭定时器", id)
			} else {
				global.Logger.Infof("用户:%d 主动退出匹配队列，关闭定时器", id)
			}
			return

		case t := <-c.limitTimer:
			// 创建定时器
			timer = time.NewTimer(time.Duration(t) * time.Second)
		case <-timer.C:
			matchedList.RemoveUserFromQueue(id)
			c.send <- newMessage(global.MatchTimeout, []byte("匹配超时"))
			return
		}
	}
}

func (c *Client) writePump(userID int64) {
	defer func() {
		c.hub.unregister <- userID
		err := c.conn.Close()
		if err != nil {
			global.Logger.Error("关闭连接时发生错误:", err)
			return
		}
		global.Logger.Infof("已与用户:%d 断开连接", userID)
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					global.Logger.Error(err)
				}
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message.message)
			if err != nil {
				global.Logger.Error(err)
				return
			}
			if message.messageType == global.MatchExit {
				c.close <- false
			} else if message.messageType != global.MatchMsg {
				c.close <- true
				return
			}
		}
	}
}
