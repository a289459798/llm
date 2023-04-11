package model

import (
	"errors"
	gogpt "github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"time"
)

type Record struct {
	ID          uint32    `gorm:"primary_key" json:"id"`
	Uid         uint32    `json:"uid" gorm:"index:idx_uid"`
	Type        string    `gorm:"type:varchar(20)" json:"type"`
	ChatId      string    `gorm:"type:varchar(50)" json:"chat_id"`
	Title       string    `gorm:"type:varchar(50)" json:"title"`
	Content     string    `json:"content" gorm:"type:text"`
	ShowContent string    `json:"show_content" gorm:"type:text"`
	Result      string    `json:"result" gorm:"type:text"`
	ShowResult  string    `json:"show_result" gorm:"type:text"`
	Model       string    `json:"model" gorm:"type:varchar(20)"`
	Platform    string    `json:"platform" gorm:"type:varchar(20)"`
	IsDelete    bool      `json:"is_delete" gorm:"default:0"`
	CreatedAt   time.Time `gorm:"column:created_at;index:idx_created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt    time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
}

func (r Record) GetMessage(db *gorm.DB, user AIUser) ([]gogpt.ChatCompletionMessage, bool, error) {
	var message []gogpt.ChatCompletionMessage
	prompt, err := getContent(r.Type)
	if err != nil {
		return nil, false, err
	}
	message = append(message, gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleUser,
		Content: "不要回复与下面问题无关的问题",
	}, gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleUser,
		Content: prompt,
	}, gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleAssistant,
		Content: "好的",
	})

	isFirst := true
	if r.ChatId != "" {
		maxToken := 300
		strLen := 100
		if user.IsVip() {
			maxToken = 800
			strLen = 200
		}
		var records []Record
		db.Raw("select id, content, LEFT(result, ?) as result from gpt_record where uid = ? and chat_id = ? and is_delete = 0 order by id desc limit ?", strLen, r.Uid, r.ChatId, 10).Scan(&records)
		if len(records) > 0 {
			isFirst = false
			totalLen := 0
			for i := len(records) - 1; i >= 0; i-- {
				message = append(message, gogpt.ChatCompletionMessage{
					Role:    "user",
					Content: records[i].Content,
				})
				message = append(message, gogpt.ChatCompletionMessage{
					Role:    "assistant",
					Content: records[i].Result,
				})

				totalLen += len([]rune(records[i].Content)) + len([]rune(records[i].Result))
				if totalLen > maxToken {
					break
				}
			}
		}
	}

	return message, isFirst, nil
}

func getContent(t string) (string, error) {
	prompt := map[string]string{
		"report/week": "请帮我把以下的工作内容填充为一篇完整的周报包含本周内容、下周计划、本周总结,用 markdown 格式以分点叙述的形式输出",
	}
	if s, ok := prompt[t]; ok {
		return s, nil
	}
	return "", errors.New("类型不存在")
}
