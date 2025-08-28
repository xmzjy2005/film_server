package middleware

import (
	"errors"
	"film_server/config"
	"film_server/model/system"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
)

// 权限验证 token中间件验证
func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求头中获取token
		authToken := c.Request.Header.Get("auth-token")
		//如果没有登录信息则清退
		if authToken == "" {
			system.CustomResult(http.StatusUnauthorized, system.SUCCESS, nil, "用户未授权，请登录", c)
			c.Abort()
			return
		}
		//解析token的信息
		uc, err := system.ParseToken(authToken)
		if err != nil {
			log.Println(err, "解析token错误")
			c.Abort()
			return
		}
		if uc == nil {
			log.Println("没有解析出数据")
			c.Abort()
			return
		}
		//从redis中获取token是否存在，存在则刷新
		t := system.GetUserTokenById(uc.UserID)
		//如果redis里面的 token为空
		if len(t) <= 0 {
			system.CustomResult(http.StatusUnauthorized, system.SUCCESS, nil, "身份证信息已失效,redis没有记录", c)
			c.Abort()
			return
		}
		if t != authToken {
			system.CustomResult(http.StatusUnauthorized, system.SUCCESS, nil, "账号在其他设备登录", c)
			c.Abort()
			return
		}
		//如果redis一样，只是过期了
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) && t == authToken {
				newToken, _ := system.GetToken(uc.UserID, uc.UserName)
				//将token存入redis
				_ = system.SaveUserToken(newToken, uc.UserID)
				//返回新的token
				c.Header("new-token", newToken)
				//解析新的token
				uc, _ = system.ParseToken(newToken)
			} else {
				system.CustomResult(http.StatusUnauthorized, system.SUCCESS, nil, err.Error(), c)
				c.Abort()
				return
			}
		}
		//将uc放至gin的上下文中
		c.Set(config.AuthUserClaims, uc)
		c.Next()
	}
}
