package utils

import (
	"chatgpt-tools/model"
	"gorm.io/gorm"
)

func GetSuanLi(t string, uid uint32, db *gorm.DB) int {
	suanli := 1
	switch t {
	case "image/create":
	case "image/edit":
	case "image/createMulti":
		suanli = 5
		break
	case "creation/article":
	case "code/generate":
	case "code/exam":
		suanli = 2
		break
	case "chat/chat":
		user := &model.User{}
		db.Joins("inner join gpt_ai on gpt_user.id=gpt_ai.uid").Where("gpt_user.id = ?", uid).Where("gpt_ai.status = 1").Find(&user)
		if user.IsVip() {
			suanli = 2
		}
		break

	}
	return suanli
}
