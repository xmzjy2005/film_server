package logic

import (
	"errors"
	"film_server/model/system"
	"film_server/plugin/common/util"
)

type UserLogic struct{}

var UL *UserLogic

// 用户登录
func (ul *UserLogic) UserLogin(account, password string) (token string, err error) {
	//根据账号密码获取用户信息
	u := system.GetUserByNameOrEmail(account)
	if u == nil {
		return "", errors.New("用户信息不存在")
	}
	//校验用户信息
	encrypt := util.PasswordEncrypt(password, u.Salt)
	if encrypt != u.Password {
		return "", errors.New("密码不正确")
	}
	//获取token，根据jwt形式，加密格式，并且用私钥签名返回token
	token, err = system.GetToken(u.ID, u.UserName)
	err = system.SaveUserToken(token, u.ID)
	return token, nil
}

// 获取用户信息
func (ul *UserLogic) GetUerInfo(id uint) system.UserInfoVo {
	u := system.GetUserById(id)
	vo := system.UserInfoVo{
		Id:       u.ID,
		UserName: u.UserName,
		Email:    u.Email,
		Gender:   u.Gender,
		NickName: u.NickName,
		Avatar:   u.Avatar,
		Status:   u.Status,
	}
	return vo
}
