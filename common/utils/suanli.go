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
		suanli = 10
		break
	case "creation/article":
	case "code/generate":
	case "code/exam":
		suanli = 2
		break
	case "chat/chat":
		ai := &model.AI{}
		db.Where("uid = ?", uid).Where("status = 1").Find(&ai)
		if ai.ID > 0 {
			suanli = 2
		}
		break

	}
	return suanli
}
