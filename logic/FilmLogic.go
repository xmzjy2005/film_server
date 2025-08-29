package logic

import "film_server/model/system"

type FilmLogic struct {
}

var FL *FilmLogic

// 获取分类树
func (fl *FilmLogic) GetFilmClassTree() system.CategoryTree {
	return system.GetCategoryTree()
}
