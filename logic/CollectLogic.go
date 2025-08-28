package logic

import "film_server/model/system"

type CollectLogic struct {
}

var CollectL *CollectLogic

// 获取采集站列表
func (cl *CollectLogic) GetFilmSourceList() []system.FilmSource {
	return system.GetCollectSourceList()
}
