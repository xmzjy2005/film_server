package controller

import (
	"film_server/logic"
	"film_server/model/system"
	"github.com/gin-gonic/gin"
)

func FilmSourceList(c *gin.Context) {
	list := logic.CollectL.GetFilmSourceList()
	system.Success(list, "资源站点列表获取成功", c)
	return
}
