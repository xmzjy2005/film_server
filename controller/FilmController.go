package controller

import (
	"film_server/logic"
	"film_server/model/system"
	"github.com/gin-gonic/gin"
)

//--------------影视分类处理---------------

// 返回影视分类树
func FilmClassTree(c *gin.Context) {
	tree := logic.FL.GetFilmClassTree()
	system.Success(tree, "影视分类信息获取成功", c)
	return
}
