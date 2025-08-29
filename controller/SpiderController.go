package controller

import (
	"film_server/logic"
	"film_server/model/system"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 开启执行任务
func StartSpider(c *gin.Context) {
	var cp system.CollectParams
	err := c.ShouldBindJSON(&cp)
	if err != nil {
		system.Failed("请求参数异常", c)
		return
	}
	fmt.Println(cp)
	if cp.Time == 0 {
		system.Failed("采集时长不能为0", c)
	}
	if cp.Batch {
		//批量
		if cp.Ids == nil || len(cp.Ids) <= 0 {
			system.Failed("Ids为空，请选择批量采集记录", c)
			return
		}
		//todo 执行批量采集
	} else {
		//单个采集
		if len(cp.Id) <= 0 {
			system.Failed("Id为空", c)
			return
		}
		logic.SL.StartCollect(cp.Id, cp.Time)

	}
	// 返回成功执行的信息
	system.SuccessOnlyMsg("采集任务已成功开启!!!", c)
}
