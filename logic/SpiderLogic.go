package logic

import (
	"errors"
	"film_server/model/system"
	"film_server/plugin/spider"
	"log"
)

type SpiderLogic struct{}

var SL *SpiderLogic

// 执行对指定站点的采集任务
func (sl *SpiderLogic) StartCollect(id string, h int) error {
	//判断id是否存在redis中
	fs := system.FindCollectSourceById(id)
	if fs == nil {
		return errors.New("id错误，没有找到采集站点信息")
	}
	go func() {
		err := spider.HandleCollect(id, h)
		if err != nil {
			log.Printf("采集任务执行失败,%s,%s", id, err)
		}
	}()

	return nil
}
