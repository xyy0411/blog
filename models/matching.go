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
	ID       int64  `gorm:"primaryKey" json:"id"`
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	ExpireAt int64  `grom:"default:2000" json:"expire_at"`

	BlockUsers      []BlockUser      `json:"block_users"`
	OnlineSoftwares []OnlineSoftware `json:"online_softwares"`
}

type BlockUser struct {
	ID         int64 `gorm:"primaryKey" json:"id"`
	MatchingID int64 `gorm:"uniqueIndex:idx_match_block;constraint:OnDelete:CASCADE" json:"matching_id"`
	UserID     int64 `gorm:"uniqueIndex:idx_match_block" json:"user_id"`
}

func (BlockUser) TableName() string {
	return "block_user"
}

type OnlineSoftware struct {
	ID         int64  `gorm:"primaryKey" json:"id"`
	MatchingID int64  `gorm:"uniqueIndex:idx_match_soft" json:"matching_id"`
	Name       string `gorm:"uniqueIndex:idx_match_soft" json:"name"`
	Type       int8   `json:"type"`
}

func (OnlineSoftware) TableName() string {
	return "online_software"
}

type MatchingSoftWare struct {
	Name string `json:"name"`
	Type int8   `json:"type"` // 0 主副皆可 1 仅主 2 仅副
}

type MatchingSoftWares []MatchingSoftWare

// IsMatch 检查是否匹配
func (user1 Matching) IsMatch(user2 Matching) (isM bool) {
	for _, blockUser := range user1.BlockUsers {
		if blockUser.UserID == user2.UserID {
			return
		}
	}
	for _, blockUser := range user2.BlockUsers {
		if blockUser.UserID == user1.UserID {
			return
		}
	}

	// 判断可用的软件
	var matchInfo MatchingSoftWares
	matchInfoMap := make(map[string]MatchingSoftWare) // 用于去重
	for _, i2 := range user1.OnlineSoftwares {
		for _, s := range user2.OnlineSoftwares {
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
	for _, s := range user1.OnlineSoftwares {
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
