package spider

import (
	"encoding/json"
	"errors"
	"film_server/model/collect"
	"film_server/model/system"
	"film_server/plugin/common/conver"
	"film_server/plugin/util"
	"fmt"
)

// jsonCollect 处理返回值为json格式的采集数据
type jsonCollect struct {
}

func (jc *jsonCollect) GetCategoryTree(r util.RequestInfo) (*system.CategoryTree, error) {
	//设置请求GET参数
	r.Params.Set("ac", "list")
	r.Params.Set("pg", "1")
	//执行请求
	util.ApiGet(&r)
	if len(r.Resp) <= 0 {
		return nil, errors.New("爬虫数据为空")
	}
	filmListPage := collect.FilmListPage{}
	//讲字节数据的json字符串，转换为结构体对象
	fmt.Println(string(r.Resp), "responce data")
	err := json.Unmarshal(r.Resp, &filmListPage)
	if err != nil {
		return nil, err
	}
	//转为层级树
	cl := filmListPage.Class
	tree := conver.GenCategoryTree(cl)
	//存储到redis
	_ = collect.SaveFilmClass(cl)
	return tree, nil
}

// 获取分页数量
func (jc *jsonCollect) GetPageCount(r util.RequestInfo) (int, error) {
	if len(r.Params.Get("ac")) <= 0 {
		r.Params.Set("ac", "detail")
	}
	r.Params.Set("pg", "1")
	util.ApiGet(&r)
	if len(r.Resp) <= 0 {
		return 0, errors.New("responese is empty")
	}

	res := collect.CommonPage{}
	err := json.Unmarshal(r.Resp, &res)
	if err != nil {
		return 0, err
	}
	count := int(res.PageCount)
	return count, nil
}

// 获取分页列表
func (jc *jsonCollect) GetFilmDetail(r util.RequestInfo) (list []system.MovieDetail, err error) {
	//设置分页请求
	r.Params.Set("ac", "detail")
	util.ApiGet(&r)
	if len(r.Resp) <= 0 {
		return
	}
	//影视详情
	detailPage := collect.FilmDetailLPage{}

	//序列化
	err = json.Unmarshal(r.Resp, &detailPage)
	if err != nil {
		return
	}
	//处理详情信息
	list = conver.ConvertFilmDetails(detailPage.List)
	return
}
