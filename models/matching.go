package models

type MatchEvent struct {
	Type    string `json:"type"` // queueing, matched, cancelled, error
	SelfID  int64  `json:"self_id"`
	PeerID  int64  `json:"peer_id"`
	Message string `json:"message"`
	MatchID string `json:"match_id"`
	Code    int    `json:"code,omitempty"`
}

type Matching struct {
	ID             uint              `gorm:"primary_key"`
	UserID         int64             `json:"user_id"`
	UserName       string            `json:"user_name"`
	GroupID        int64             `json:"group_id"`
	LimitTime      int64             `json:"limit_time"`
	BlockUser      []*BlockUser      `gorm:"foreignKey:MatchingID" json:"block_user"`
	OnlineSoftware []*OnlineSoftware `gorm:"foreignKey:MatchingID" json:"online_software"`
}

type BlockUser struct {
	ID         int64 `gorm:"primary_key" json:"id"`
	MatchingID int64 `json:"matching_id"`
	UserID     int64 `gorm:"unique" json:"user_id"`
}

type OnlineSoftware struct {
	ID         int64  `gorm:"primary_key" json:"id"`
	MatchingID int64  `json:"matching_id"`
	Name       string `gorm:"unique" json:"name"`
	Type       int8   `gorm:"unique" json:"type"` // 0 主副皆可 1 仅主 2 仅副
}

type SendMatchMessage struct {
	UserID   int64             `json:"user_id"`
	Software MatchingSoftWares `json:"software"`
}

type MatchingSoftWare struct {
	Name string `json:"name"`
	Type int8   `json:"type"` // 0 主副皆可 1 仅主 2 仅副
}

type MatchingSoftWares []MatchingSoftWare

// IsMatching 检查是否匹配
func (user1 Matching) IsMatching(user2 Matching) (isM bool) {
	blocked := make(map[int64]struct{})
	for _, blockUser := range user1.BlockUser {
		blocked[blockUser.UserID] = struct{}{}
	}
	for _, blockUser := range user2.BlockUser {
		if _, exists := blocked[blockUser.UserID]; exists {
			return
		}
	}

	// 判断可用的软件
	var matchInfo MatchingSoftWares
	matchInfoMap := make(map[string]MatchingSoftWare) // 用于去重
	for _, i2 := range user1.OnlineSoftware {
		for _, s := range user2.OnlineSoftware {
			/*			// 最优先匹配软件相同 类型为互为12
						if s.Name == i2.Name && math.Abs(float64(s.Type-i2.Type)) == 1 && s.Type != 0 && i2.Type != 0 {
							global.Logger.Info("匹配成功")
							return true
						}*/
			if s.Type-i2.Type == 0 && s.Type+i2.Type != 0 {
				continue
			}
			// 检查软件是否已经存在于 matchInfo 中
			if _, exist := matchInfoMap[s.Name]; !exist {
				matchInfoMap[s.Name] = MatchingSoftWare{
					Name: s.Name,
					Type: s.Type,
				}
				matchInfo = append(matchInfo, matchInfoMap[s.Name])
			}
		}
	}

	// 匹配的软件至少要有一个相同的
	exists := false
	matchInfoMapForExistCheck := make(map[string]struct{})
	for _, ware := range matchInfo {
		matchInfoMapForExistCheck[ware.Name] = struct{}{}
	}
	// 检查 user1 的软件是否在 matchInfo 中
	for _, s := range user1.OnlineSoftware {
		if _, exist := matchInfoMapForExistCheck[s.Name]; exist {
			exists = true
			break
		}
	}
	// 如果匹配到至少一个相同的软件，返回 true
	if len(matchInfo) > 0 && exists {
		isM = true
	}

	return
}
