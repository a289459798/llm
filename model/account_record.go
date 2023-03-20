package model

import (
	"time"
)

type AccountRecord struct {
	ID            uint32    `gorm:"primary_key" json:"id"`
	Uid           uint32    `json:"uid" gorm:"index:idx_uid"`
	RecordId      uint32    `json:"record_id"`
	Way           uint8     `json:"type" gorm:"COMMENT:方式"`
	Type          string    `gorm:"type:varchar(20)" json:"type"`
	Amount        uint32    `json:"amount"`
	CurrentAmount uint32    `json:"current_amount"`
	CreatedAt     time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt      time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (a AccountRecord) GetType() string {
	types := map[string]string{
		"open":              "活跃",
		"convert/translate": "翻译",
		"share":             "分享",
		"share_follow":      "分享打开",
		"chat/chat":         "对话",
		"ad":                "看广告",
		"creation/article":  "作文",
		"image/createMulti": "画画",
		"code/generate":     "代码生成",
		"report/week":       "周报",
		"report/day":        "日报",
		"welfare":           "福利",
		"image/pic2pic":     "以图绘图",
		"chat/introduce":    "自我介绍",
	}
	if v, ok := types[a.Type]; ok {
		return v
	}
	return "其他根工具"
}
