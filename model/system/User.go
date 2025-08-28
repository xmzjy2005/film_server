package system

import (
	"film_server/config"
	"film_server/plugin/db"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	UserName string `json:"userName"` // 用户名
	Password string `json:"password"` // 密码
	Salt     string `json:"salt"`     // 盐值
	Email    string `json:"email"`    // 邮箱
	Gender   int    `json:"gender"`   // 性别
	NickName string `json:"nickName"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
	Status   int    `json:"status"`   // 状态
	Reserve1 string `json:"reserve1"` // 预留字段 3
	Reserve2 string `json:"reserve2"` // 预留字段 2
	Reserve3 string `json:"reserve3"` // 预留字段 1
	//LastLongTime time.Time `json:"LastLongTime"` // 最后登录时间
}

// 设置user的表名
func (u *User) TableName() string { return config.UserTableName }

func GetUserByNameOrEmail(userName string) *User {
	var u *User
	if err := db.Mdb.Where("user_name = ? or email = ? ", userName, userName).First(&u).Error; err != nil {
		log.Println(err)
		return nil
	}
	return u
}

// 通过iD获取用户信息
func GetUserById(id uint) User {
	user := User{Model: gorm.Model{ID: id}}
	db.Mdb.First(&user)
	return user
}
