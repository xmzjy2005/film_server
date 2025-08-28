package system

import (
	"encoding/json"
	"film_server/config"
	"film_server/plugin/db"
	"log"
)

// 分类结构
type Category struct {
	Id   int64  `json:"id"`   //分类id
	Pid  int64  `json:"pid"`  //父级分类
	Name string `json:"name"` //分类名称
	Show bool   `json:"show"` //是否展示
}
type CategoryTree struct {
	*Category
	Children []*CategoryTree `json:"children"` //子分类信息
}

// 查询分类信息是否存在
func ExistsCategoryTree() bool {
	exists, err := db.Rdb.Exists(db.Cxt, config.CategoryTreeKey).Result()
	if err != nil {
		log.Println("ExistsCategoryTree Error", err)
	}
	return exists == 1
}

// 保存分类
func SaveCategoryTree(tree *CategoryTree) error {
	data, _ := json.Marshal(tree)
	return db.Rdb.Set(db.Cxt, config.CategoryTreeKey, data, config.FilmExpired).Err()
}

// 从redis获取树
func GetCategoryTree() CategoryTree {
	data := db.Rdb.Get(db.Cxt, config.CategoryTreeKey).Val()
	tree := CategoryTree{}
	_ = json.Unmarshal([]byte(data), &tree)
	return tree
}

// 获取子集分类
func GetChildrenTree(id int64) []*CategoryTree {
	tree := GetCategoryTree()
	for _, t := range tree.Children {
		if t.Id == id {
			return t.Children
		}
	}
	return nil
}
