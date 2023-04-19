package model

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type AIUser struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Uid       uint32    `gorm:"after:id;index:ik_uid" json:"uid"`
	OpenId    string    `gorm:"type:varchar(64);uniqueIndex:uk_openid" json:"open_id"`
	UnionId   string    `gorm:"type:varchar(64)"json:"union_id"`
	Subscribe bool      `json:"subscribe"`
	JoinGroup bool      `json:"join_group" gorm:"default:0"`
	AppKey    string    `json:"app_key" gorm:"type:varchar(32)"`
	Channel   string    `json:"channel" gorm:"type:varchar(32)"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_at,omitempty"`
	UpdateAt  time.Time `gorm:"column:update_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP  on update current_timestamp" json:"update_at,omitempty"`
	Vip       AIUserVip `gorm:"foreignKey:uid;references:uid"`
}

type UserLogin struct {
	OpenID       string
	UnionID      string
	Channel      string
	AppKey       string
	AccessExpire int64
	AccessSecret string
}

func (user AIUser) Login(db *gorm.DB, userLogin UserLogin) (*AIUser, string, error) {
	db.Where("open_id = ?", userLogin.OpenID).First(&user)
	if user.Uid == 0 {
		// 判断UnionID是否存在
		db.Where("union_id = ?", userLogin.UnionID).First(&user)
		if user.Uid == 0 || userLogin.UnionID == "" {
			tx := db.Begin()
			// 创建用户
			u := &User{}
			tx.Create(u)
			user.OpenId = userLogin.OpenID
			user.UnionId = userLogin.UnionID
			user.AppKey = userLogin.AppKey
			user.Channel = userLogin.Channel
			user.Uid = u.ID
			user.ID = 0
			err := tx.Create(&user).Error
			if err != nil {
				tx.Rollback()
				return nil, "", err
			}
			if userLogin.Channel != "" {
				distributorUid, err := strconv.Atoi(userLogin.Channel)
				if err == nil {
					DistributorRecord{
						DistributorUid: uint32(distributorUid),
						Uid:            user.Uid,
					}.Create(db)
				}
			}
			tx.Commit()
		} else {
			newUser := &AIUser{}
			newUser.OpenId = userLogin.OpenID
			newUser.UnionId = userLogin.UnionID
			newUser.AppKey = userLogin.AppKey
			newUser.Channel = userLogin.Channel
			newUser.Uid = user.Uid
			err := db.Create(&newUser).Error
			if err != nil {
				return nil, "", err
			}
		}
	} else if user.UnionId == "" {
		user.UnionId = userLogin.UnionID
		db.Save(user)
	}
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Unix() + userLogin.AccessExpire
	claims["iat"] = time.Now().Unix()
	claims["uid"] = user.Uid
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(userLogin.AccessSecret))
	if err != nil {
		return nil, "", err
	}
	return &user, tokenString, nil
}

func (user AIUser) Find(db *gorm.DB) AIUser {
	db.Where("uid = ?", user.Uid).Preload("Vip").First(&user)
	return user
}

func (user AIUser) IsVip() bool {
	return time.Now().Unix() < user.Vip.VipExpiry.Unix()
}

func (user AIUser) IsJoinGroup(db *gorm.DB) bool {
	var c int64
	db.Model(AIUser{}).Where("uid = ?", user.Uid).Where("join_group = 1").Count(&c)
	return c > 0
}

func (user AIUser) SetVip(db *gorm.DB, vipCode *VipCode) error {
	if user.Uid == 0 {
		return errors.New("用户不存在")
	}
	return db.Transaction(func(tx *gorm.DB) error {
		var err error
		// 设置VIP过期时间
		if user.IsVip() {
			user.Vip.VipExpiry = user.Vip.VipExpiry.AddDate(0, 0, int(vipCode.Day))
		} else {
			user.Vip.VipExpiry, _ = time.ParseInLocation("2006-01-02 15:04:05", time.Now().AddDate(0, 0, int(vipCode.Day)).Format("2006-01-02")+" 23:59:59", time.Local)
		}
		if user.Vip.ID > 0 {
			err = tx.Save(&user.Vip).Error
		} else {
			user.Vip.Uid = user.Uid
			user.Vip.VipId = vipCode.VipId
			user.Vip.Amount = vipCode.Vip.Amount
			err = tx.Create(&user.Vip).Error
		}
		if err != nil {
			return err
		}

		// 增加VIP算力
		var vipAmount = vipCode.Vip.Amount
		if vipAmount > 0 {
			amount := NewAccount(tx).GetAccount(user.Uid, time.Now())
			amount.ChatAmount += vipAmount
			err = tx.Save(&amount).Error
			if err != nil {
				return err
			}

			// 插入算力明细
			err := tx.Create(&AccountRecord{
				Uid:           user.Uid,
				RecordId:      0,
				Way:           1,
				Type:          "vip",
				Amount:        vipAmount,
				CurrentAmount: amount.Amount + vipAmount,
			}).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}
