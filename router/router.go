package router

import (
	"film_server/controller"
	"film_server/plugin/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	//开启跨域
	//r.Use()
	r.POST("/login", controller.Login)
	//后台登录首页
	manageRoute := r.Group("/manage")
	//这边怎么取的token？login成功后，是放到new-token里面的头部，然后呢？
	manageRoute.Use(middleware.AuthToken())
	{
		manageRoute.GET("/index", controller.ManageIndex)

		//用户相关
		userRoute := manageRoute.Group("user")
		{
			userRoute.GET("/info", controller.UserInfo)
		}

		//采集相关
		collect := manageRoute.Group("/collect")
		{
			collect.GET("/list", controller.FilmSourceList)
		}

		//数据采集
		spiderRoute := manageRoute.Group("/spider")
		{
			spiderRoute.POST("/start", controller.StartSpider)
		}

		//影视管理
		filmRoute := manageRoute.Group("film")
		{
			filmRoute.GET("/class/tree", controller.FilmClassTree)
		}
	}
	return r
}
