package controller

import (
	"film_server/model/system"
	"github.com/gin-gonic/gin"
)

func ManageIndex(c *gin.Context) {
	system.SuccessOnlyMsg("后台中心", c)
	return
}
