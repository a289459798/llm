package utils

import (
	"chatgpt-tools/model"
	"encoding/json"
	"gorm.io/gorm"
)

func GetSuanLi(uid uint32, t string, params string, db *gorm.DB) int {
	suanli := 1
	paramsMap := make(map[string]interface{})
	if params != "" {
		json.Unmarshal([]byte(params), &paramsMap)

	}
	switch t {
	case "image/create", "image/edit", "image/createMulti":
		suanli = 5
		if clarity, ok := paramsMap["clarity"]; ok {
			if clarity == "superhigh" {
				suanli = 20
			} else if clarity == "high" {
				suanli = 10
			}
		}
		if number, ok := paramsMap["number"]; ok {
			suanli = suanli * int(number.(float64))
		}
		break
	case "image/img2text":
		suanli = 5
		break
	case "image/pic-repair":
		suanli = 10
		break
	case "creation/article", "code/generate", "code/exam":
		suanli = 2
		break
	case "chat/chat":
		user := &model.AIUser{}
		db.Joins("inner join gpt_ai on gpt_ai_user.uid=gpt_ai.uid").Where("gpt_ai_user.uid = ?", uid).Where("gpt_ai.status = 1").Find(&user)
		if user.IsVip() {
			suanli = 2
		}
		break

	}
	return suanli
}
