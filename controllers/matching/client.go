package matching

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/models"
	matchingrepo "github.com/xyy0411/blog/repositories/matching"
)

type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan []byte // 发送消息的通道
	limitTimer  chan int64  // 限制匹配时间的通道
	close       chan bool   // 关闭通道
	connectedAt time.Time   // 连接时间，用于计算 Duration
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:         hub,
		conn:        conn,
		send:        make(chan []byte, 256),
		limitTimer:  make(chan int64),
		close:       make(chan bool),
		connectedAt: time.Now(),
	}
}

func (c *Client) checkLimitTimer(id int64) {
	timer := time.NewTimer(24 * time.Hour)
	timer.Stop()
	defer func() {
		global.Logger.Infof("用户:%d 已退出匹配队列,关闭定时器", id)
		if timer != nil {
			timer.Stop()
		}
	}()

	for {
		select {
		case <-c.close:
			return
		case t := <-c.limitTimer:
			// 创建或重置定时器
			timer = time.NewTimer(time.Duration(t) * time.Second)
		case <-timer.C:
			matchedList.RemoveUserFromQueue(id)

			// 获取用户信息用于保存记录
			repo := matchingrepo.NewRepo(global.DB)
			user, err := repo.GetByUserID(id)
			if err != nil {
				global.Logger.Errorf("获取用户信息失败: %v", err)
				user = models.Matching{UserID: id, UserName: ""}
			}

			// 计算匹配持续时间（秒）
			duration := int(time.Now().Sub(c.connectedAt).Seconds())

			// 保存匹配申请记录（匹配超时）
			matchedList.saveMatchingApplication(id, user.UserName, false, duration, "", models.ExitReasonTimeout)

			event := models.MatchEvent{
				Type:    "error",
				SelfID:  id,
				PeerID:  0,
				Message: "匹配超时,已退出匹配队列",
				Code:    http.StatusRequestTimeout,
			}
			msg, _ := json.Marshal(event)
			c.send <- msg
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
		c.close <- true
		global.Logger.Debugf("已与用户:%d 断开连接", userID)
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
