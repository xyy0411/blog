package matching

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xyy0411/blog/global"
	"github.com/xyy0411/blog/models"
	matchingrepo "github.com/xyy0411/blog/repositories/matching"
	"github.com/xyy0411/blog/resp"
	"gorm.io/gorm"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	matchedList = NewMatchingManager()
	lock        sync.RWMutex
)

func getRepo() *matchingrepo.Repo {
	return matchingrepo.NewRepo(global.DB)
}

func parseUserIDParam(ctx *gin.Context) (int64, bool) {
	id := ctx.Param("user_id")
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		resp.Error(ctx, http.StatusBadRequest, "user_id 参数无效")
		return 0, false
	}
	return userID, true
}

func sendEvent(client *Client, event models.MatchEvent) {
	msg, _ := json.Marshal(event)
	client.send <- msg
}

func syncQueueUser(userID int64) {
	_, inQueue := matchedList.matchedList.Load(userID)
	if !inQueue {
		return
	}

	user, err := getRepo().GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			matchedList.RemoveUserFromQueue(userID)
			return
		}
		global.Logger.Errorf("同步匹配队列用户失败，user_id=%d，err=%v", userID, err)
		return
	}

	matchedList.AddUserToQueue(&user)
}

func ensureDefaultMatchingProfile(userID int64) (models.Matching, error) {
	defaultSoftwares := []models.OnlineSoftware{
		{Name: "uu", Type: 0},
		{Name: "to", Type: 0},
	}

	user, err := getRepo().GetByUserID(userID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Matching{}, err
		}

		return models.Matching{
			UserID:          userID,
			UserName:        strconv.FormatInt(userID, 10),
			OnlineSoftwares: defaultSoftwares,
		}, nil
	}

	if len(user.OnlineSoftwares) == 0 {
		user.OnlineSoftwares = defaultSoftwares
	}

	return user, nil
}

func UpdateProfileName(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}
	var input struct {
		UserName string `json:"user_name"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		resp.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := getRepo().UpdateName(userID, input.UserName); err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "更新用户名失败")
		return
	}
	resp.OK(ctx, "更新用户名成功！", nil)
}

func CreateMatchingProfile(ctx *gin.Context) {
	var input struct {
		UserID          int64                   `json:"user_id"`
		UserName        string                  `json:"user_name"`
		ExpireAt        int64                   `json:"expire_at"`
		OnlineSoftwares []models.OnlineSoftware `json:"online_softwares"`
		BlockUserIDs    []int64                 `json:"block_user_ids"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		resp.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if input.UserID == 0 {
		resp.Error(ctx, http.StatusBadRequest, "user_id 不能为空")
		return
	}

	if strings.TrimSpace(input.UserName) == "" {
		input.UserName = strconv.FormatInt(input.UserID, 10)
	}

	repo := getRepo()
	info, err := repo.GetByUserID(input.UserID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配资料失败")
		return
	}

	if info.UserName == input.UserName {
		resp.OK(ctx, "用户信息已是最新无需更新", nil)
		return
	}

	match := models.Matching{
		UserID:          input.UserID,
		UserName:        input.UserName,
		OnlineSoftwares: input.OnlineSoftwares,
	}
	if input.ExpireAt > 0 {
		match.ExpireAt = input.ExpireAt
	}

	if len(input.BlockUserIDs) > 0 {
		match.BlockUsers = make([]models.BlockUser, 0, len(input.BlockUserIDs))
		for _, blocked := range input.BlockUserIDs {
			if blocked == 0 || blocked == input.UserID {
				continue
			}
			match.BlockUsers = append(match.BlockUsers, models.BlockUser{UserID: blocked})
		}
	}

	if err := repo.CreateMatchingWithChildren(&match); err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "创建用户失败")
		return
	}

	resp.OK(ctx, "创建成功", map[string]any{"matching": match})
}

func GetMatchingProfile(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	match, err := getRepo().GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error(ctx, http.StatusNotFound, "未找到匹配资料")
			return
		}
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配资料失败")
		return
	}

	resp.OK(ctx, "", map[string]any{"matching": match})
}

func GetMatchingSoftwareList(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	match, err := getRepo().GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error(ctx, http.StatusNotFound, "未找到匹配资料")
			return
		}
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配资料失败")
		return
	}

	resp.OK(ctx, "", map[string]any{"online_softwares": match.OnlineSoftwares})
}

func GetMatchingBlockUserList(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	match, err := getRepo().GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error(ctx, http.StatusNotFound, "未找到匹配资料")
			return
		}
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配资料失败")
		return
	}

	blockUserIDs := make([]int64, 0, len(match.BlockUsers))
	for _, blockUser := range match.BlockUsers {
		blockUserIDs = append(blockUserIDs, blockUser.UserID)
	}

	resp.OK(ctx, "", map[string]any{"block_user_ids": blockUserIDs})
}

func UpdateMatchingExpire(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	var input struct {
		ExpireAt int64 `json:"expire_at"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		resp.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if input.ExpireAt <= 0 {
		resp.Error(ctx, http.StatusBadRequest, "expire_at 必须大于 0")
		return
	}

	if err := getRepo().UpdateExpire(userID, input.ExpireAt); err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "更新 expire_at 失败")
		return
	}

	syncQueueUser(userID)
	resp.OK(ctx, "更新成功", nil)
}

func GetMatchingExpire(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	match, err := getRepo().GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error(ctx, http.StatusNotFound, "未找到匹配资料")
			return
		}
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配资料失败")
		return
	}

	resp.OK(ctx, "", map[string]any{"expire_at": match.ExpireAt})
}

func AddMatchingSoftware(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	var input struct {
		Name string `json:"name"`
		Type int8   `json:"type"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		resp.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		resp.Error(ctx, http.StatusBadRequest, "软件名称不能为空")
		return
	}
	if input.Type < 0 || input.Type > 2 {
		resp.Error(ctx, http.StatusBadRequest, "软件类型必须为 0/1/2")
		return
	}

	repo := getRepo()
	match, err := repo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error(ctx, http.StatusNotFound, "未找到匹配资料")
			return
		}
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配资料失败")
		return
	}

	for _, s := range match.OnlineSoftwares {
		if s.Name == input.Name {
			if err := repo.UpdateOnlineSoftwareType(match.ID, input.Name, input.Type); err != nil {
				resp.Error(ctx, http.StatusInternalServerError, "更新软件失败")
				return
			}
			syncQueueUser(userID)
			resp.OK(ctx, "软件更新成功", nil)
			return
		}
	}

	if err := repo.AddOnlineSoftware(match.ID, models.OnlineSoftware{Name: input.Name, Type: input.Type}); err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "添加软件失败")
		return
	}

	syncQueueUser(userID)
	resp.OK(ctx, "软件添加成功", nil)
}

func RemoveMatchingSoftware(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	var input struct {
		SoftwareName string `json:"software_name"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		resp.Error(ctx, http.StatusBadRequest, errors.Join(errors.New("software字段错误"), err))
		return
	}

	repo := getRepo()
	match, err := repo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error(ctx, http.StatusNotFound, "未找到匹配软件")
			return
		}
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配资料失败")
		return
	}

	if err := repo.RemoveOnlineSoftware(match.ID, input.SoftwareName); err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "移除软件失败")
		return
	}

	syncQueueUser(userID)
	resp.OK(ctx, "软件移除成功", nil)
}

func AddMatchingBlockUser(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	var input struct {
		TargetUserID int64 `json:"target_user_id"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		resp.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if input.TargetUserID == 0 || input.TargetUserID == userID {
		resp.Error(ctx, http.StatusBadRequest, "target_user_id 参数无效")
		return
	}

	repo := getRepo()
	match, err := repo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error(ctx, http.StatusNotFound, "未找到匹配资料")
			return
		}
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配资料失败")
		return
	}

	if err := repo.AddBlockUser(match.ID, input.TargetUserID); err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "添加屏蔽用户失败")
		return
	}

	syncQueueUser(userID)
	resp.OK(ctx, "屏蔽用户添加成功", nil)
}

func RemoveMatchingBlockUser(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	targetIDStr := ctx.Param("target_user_id")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 64)
	if err != nil || targetID <= 0 {
		resp.Error(ctx, http.StatusBadRequest, "target_user_id 参数无效")
		return
	}

	repo := getRepo()
	match, err := repo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Error(ctx, http.StatusNotFound, "未找到匹配资料")
			return
		}
		resp.Error(ctx, http.StatusInternalServerError, "查询匹配资料失败")
		return
	}

	if err := repo.RemoveBlockUser(match.ID, targetID); err != nil {
		resp.Error(ctx, http.StatusInternalServerError, "移除屏蔽用户失败")
		return
	}

	syncQueueUser(userID)
	resp.OK(ctx, "屏蔽用户移除成功", nil)
}

func GetMatchingPerson(ctx *gin.Context) {
	resp.OK(ctx, strconv.Itoa(matchedList.Len()), nil)
}

func LookMatchingStatus(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	_, inQueue := matchedList.matchedList.Load(userID)
	if !inQueue {
		resp.Error(ctx, http.StatusNotFound, "用户不在匹配队列中")
		return
	}

	resp.OK(ctx, "用户正在匹配队列中", nil)
}

func QuitMatching(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	_, inQueue := matchedList.matchedList.Load(userID)
	if !inQueue {
		resp.Error(ctx, http.StatusNotFound, "用户不在匹配队列中")
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

	sendEvent(client, models.MatchEvent{
		Type:    "cancelled",
		SelfID:  userID,
		PeerID:  0,
		Message: "你已退出匹配队列",
		Code:    200,
	})
	resp.OK(ctx, "你已退出匹配队列", nil)
}

func HandleMatching(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		resp.Error(ctx, http.StatusInternalServerError, err.Error())
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

	user, err := ensureDefaultMatchingProfile(userID)
	if err != nil {
		global.Logger.Errorf("用户:%d 连接异常:%v", userID, err)
		event := models.MatchEvent{
			Type:    "error",
			SelfID:  userID,
			PeerID:  0,
			Message: err.Error(),
			Code:    500,
		}
		sendEvent(client.client, event)
		return
	}

	if userID != user.UserID {
		sendEvent(client.client, models.MatchEvent{Type: "error", SelfID: userID, Message: "用户ID不匹配", Code: 400})
		return
	}

	if _, ok := matchedList.matchedList.Load(user.UserID); ok {
		event := models.MatchEvent{
			Type:    "error",
			SelfID:  userID,
			PeerID:  0,
			Message: "你已在匹配队列中，请勿重复匹配",
			Code:    409,
		}
		sendEvent(client.client, event)
		return
	}

	MatchHub.match <- user
	client.client.limitTimer <- user.ExpireAt
}
