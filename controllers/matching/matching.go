package matching

import (
	"encoding/json"
	"github.com/RomiChan/syncx"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/models"
	"github.com/xyy0411/blog/resp"
	"github.com/xyy0411/blog/utils"
	"net/http"
	"strconv"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		/*		ReadBufferSize:  1024,
				WriteBufferSize: 1024,*/
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// 存储匹配的用户
	matchedList = NewMatchingManager()

	lock sync.RWMutex
	/*	// 维护用户 WebSocket 连接
		connections = syncx.Map[int64, *Client]{}*/
)

type Manager struct {
	// 存储匹配的用户
	matchedList syncx.Map[int64, *models.Matching]
	// 匹配ID计数器
	matchIDCounter int64
}

func NewMatchingManager() *Manager {
	return &Manager{
		matchedList: syncx.Map[int64, *models.Matching]{},
	}
}

func (mm *Manager) Len() int {
	var listLen int

	mm.matchedList.Range(func(key int64, value *models.Matching) bool {
		listLen++
		return true
	})

	return listLen
}

func (mm *Manager) AddUserToQueue(user *models.Matching) {
	mm.matchedList.Store(user.UserID, user)
}

func (mm *Manager) RemoveUserFromQueue(userID int64) {
	mm.matchedList.Delete(userID)
}

func (mm *Manager) notifyAndRemoveUser(id int64, user models.Matching, matchID string) {
	// 检查 MatchHub 是否为 nil
	if MatchHub == nil {
		global.Logger.Error("MatchHub 未初始化")
		return
	}
	// 检查 MatchHub.clients 是否为 nil
	if MatchHub.clients == nil {
		global.Logger.Error("MatchHub.clients 未初始化")
		return
	}
	client, ok := MatchHub.clients[id]
	if !ok || client == nil {
		global.Logger.Errorf("用户 %d 的客户端未找到", user.UserID)
		return
	}
	// 检查 client.send 是否为 nil
	if client.send == nil {
		global.Logger.Errorf("用户 %d 的客户端 send 通道未初始化", user.UserID)
		return
	}
	event := utils.FormatMatchingInfo(id, user, matchID)
	mm.RemoveUserFromQueue(user.UserID)
	// 匹配结束前要关闭定时器
	client.close <- true
	msg, _ := json.Marshal(event)
	client.send <- msg
}

func (mm *Manager) MatchUsers(user models.Matching) {

	if mm.Len() == 0 {
		mm.AddUserToQueue(&user)
		return
	}

	var targetUser models.Matching
	mm.matchedList.Range(func(key int64, value *models.Matching) bool {
		if user.UserID != value.UserID && user.IsMatching(*value) {
			targetUser = *value
			return false
		}
		return true
	})

	if targetUser.UserID == 0 {
		mm.AddUserToQueue(&user)
		return
	}

	// 生成匹配ID
	mm.matchIDCounter++
	matchID := strconv.FormatInt(mm.matchIDCounter, 10)

	// 发送消息
	mm.notifyAndRemoveUser(targetUser.UserID, user, matchID)
	mm.notifyAndRemoveUser(user.UserID, targetUser, matchID)

	global.Logger.Infof("匹配成功 用户:%d <----> 用户:%d, 匹配ID:%s", user.UserID, targetUser.UserID, matchID)

	return
}

func GetMatchingPerson(ctx *gin.Context) {
	resp.OK(ctx, strconv.Itoa(matchedList.Len()), nil)
}

func LookMatchingStatus(ctx *gin.Context) {
	id := ctx.Param("user_id")
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, err)
		return
	}
	_, ok := matchedList.matchedList.Load(userID)
	if !ok {
		resp.Error(ctx, http.StatusNotFound, "你还没有加入匹配队列")
		return
	}
	resp.OK(ctx, "你正在匹配队列中", nil)
}

func QuitMatching(ctx *gin.Context) {
	id := ctx.Param("user_id")
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, err)
		return
	}

	_, ok := matchedList.matchedList.Load(userID)
	if !ok {
		resp.Error(ctx, http.StatusNotFound, "你还没有加入匹配队列")
		return
	}
	matchedList.RemoveUserFromQueue(userID)

	lock.Lock()
	client, ok := MatchHub.clients[userID]
	lock.Unlock()
	if !ok {
		resp.Error(ctx, http.StatusNotFound, "未找到用户的 WebSocket 连接")
		return
	}
	event := models.MatchEvent{
		Type:    "cancelled",
		SelfID:  userID,
		PeerID:  0,
		Message: "你已退出匹配队列",
		Code:    200,
	}
	msg, _ := json.Marshal(event)
	client.send <- msg
}

func HandleMatching(ctx *gin.Context) {
	id := ctx.Param("user_id")
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		resp.Error(ctx, http.StatusInternalServerError, err)
		return
	}
	global.Logger.Infof("已与用户:%d 建立 WebSocket 连接", userID)

	client := &userClient{
		id:     userID,
		client: NewClient(MatchHub, conn),
	}

	// 启动写消息的协程
	go client.client.writePump(userID)
	// 注册客户端
	MatchHub.register <- client
	// 启动定时器
	go client.client.checkLimitTimer(userID)

	for {
		var user models.Matching
		// 拿一下匹配的软件
		err = conn.ReadJSON(&user)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				global.Logger.Errorf("用户:%d 连接异常:%v", userID, err)
				event := models.MatchEvent{
					Type:    "error",
					SelfID:  userID,
					PeerID:  0,
					Message: err.Error(),
					Code:    500,
				}
				msg, _ := json.Marshal(event)
				client.client.send <- msg
			}
			break
		}
		if userID != user.UserID {
			event := models.MatchEvent{
				Type:    "error",
				SelfID:  userID,
				PeerID:  0,
				Message: "用户ID不匹配",
				Code:    400,
			}
			msg, _ := json.Marshal(event)
			client.client.send <- msg
			break
		}

		if _, ok := matchedList.matchedList.Load(user.UserID); ok {
			event := models.MatchEvent{
				Type:    "error",
				SelfID:  userID,
				PeerID:  0,
				Message: "你已在匹配队列中，请勿重复匹配",
				Code:    409,
			}
			msg, _ := json.Marshal(event)
			client.client.send <- msg
			break
		}

		MatchHub.match <- user
		client.client.limitTimer <- user.LimitTime
	}
}
