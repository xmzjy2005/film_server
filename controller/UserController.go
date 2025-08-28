package controller

import (
	"film_server/config"
	"film_server/logic"
	"film_server/model/system"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	//输入账号密码，正确后存入redis
	var u system.User
	username := c.PostForm("userName")
	fmt.Println(username)
	err := c.ShouldBindJSON(&u)
	if err != nil {
		fmt.Println(err.Error(), "error")
		fmt.Println(u, "u...")
		system.Failed("获取参数异常", c)
		return
	}
	if len(u.UserName) <= 0 || len(u.Password) <= 0 {
		system.Failed("用户名或者密码为空", c)
		return
	}
	//登录验证，并且生成token后存入redis
	token, err := logic.UL.UserLogin(u.UserName, u.Password)
	if err != nil {
		system.Failed(err.Error(), c)
		return
	}
	//响应存入头部，客户端根据响应头获取token，token不放在响应体而是放在响应头，说是这样更简洁符合restfull
	c.Header("new-token", token)
	system.SuccessOnlyMsg("登录成功!!!", c)

}
func UserInfo(c *gin.Context) {
	v, ok := c.Get(config.AuthUserClaims)
	if !ok {
		system.Failed("用户信息获取失败，auth验证不通过", c)
		return
	}
	uc, ok := v.(*system.UserClaims)
	if !ok {
		system.Failed("用户信息获取失败 ，断言错误", c)
	}
	info := logic.UL.GetUerInfo(uc.UserID)
	system.Success(info, "成功返回用户信息", c)
}
