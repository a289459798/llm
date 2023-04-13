package model

import (
	"time"
)

type AccountRecord struct {
	ID            uint32    `gorm:"primary_key" json:"id"`
	Uid           uint32    `json:"uid" gorm:"uniqueIndex:uk_uid_date"`
	RecordId      uint32    `json:"record_id"`
	Way           uint8     `json:"type" gorm:"COMMENT:方式"`
	Type          string    `gorm:"type:varchar(20);uniqueIndex:uk_uid_date" json:"type"`
	Amount        uint32    `json:"amount"`
	CurrentAmount uint32    `json:"current_amount"`
	CreatedAt     time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create;uniqueIndex:uk_uid_date" json:"created_at,omitempty"`
	UpdateAt      time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (a AccountRecord) GetType() string {
	types := map[string]string{
		"open":              "活跃",
		"share":             "分享",
		"share_follow":      "分享打开",
		"ad":                "看广告",
		"welfare":           "福利",
		"group":             "加群",
		"exchange":          "兑换",
		"vip":               "会员福利",
		"convert/translate": "翻译",
		"chat/chat":         "对话",
		"creation/article":  "作文",
		"image/createMulti": "画画",
		"code/generate":     "代码生成",
		"report/week":       "周报",
		"report/day":        "日报",
		"image/pic2pic":     "以图绘图",
		"chat/introduce":    "自我介绍",
		"image/ps":          "P图",
	}
	if v, ok := types[a.Type]; ok {
		return v
	}
	return "其他工具"
}
