package collect

import (
	"encoding/json"
	"film_server/config"
	"film_server/plugin/db"
)

/*
请求数据来源 redis https://360zy.com/api.php/provide/vod/at/json?ac=list
用于生成树结构
响应数据格式
{
	"code": 1,
	"msg": "数据列表",
	"page": "1",
	"pagecount": 11,
	"limit": "20",
	"total": 209,
	"list": [{
		"vod_id": 74281,
		"type_id": 35,
		"type_id_1": 3,
		"group_id": 0,
		"vod_name": "食尚玩家2025",
		"vod_sub": "SuperTaste",
		"vod_en": "shishangwanjia2025",
	"class": [{
		"type_id": 1,
		"type_pid": 0,
		"type_name": "电影"
	}, {
		"type_id": 2,
		"type_pid": 0,
		"type_name": "连续剧"
	}, {
		"type_id": 3,
		"type_pid": 0,
		"type_name": "综艺"
	}, {
*/
// 分页结构体
type FilmListPage struct {
	Code      int         `json:"code"`      //响应状态码
	Msg       string      `json:"msg"`       //数据类型
	Page      any         `json:"page"`      //页码
	PageCount int         `json:"pagecount"` //总页数
	Limit     any         `json:"limit"`     //每页数据量
	Total     int         `json:"total"`     //总数量
	List      []FilmList  `json:"lit"`       //列表
	Class     []FilmClass `json:"class"`     //影视分类信息
}

// FilmList 影视列表单部影片信息结构体
type FilmList struct {
	VodID       int64  `json:"vod_id"`        // 影片ID
	VodName     string `json:"vod_name"`      // 影片名称
	TypeID      int64  `json:"type_id"`       // 分类ID
	TypeName    string `json:"type_name"`     // 分类名称
	VodEn       string `json:"vod_en"`        // 影片名中文拼音
	VodTime     string `json:"vod_time"`      // 更新时间
	VodRemarks  string `json:"vod_remarks"`   // 更新状态
	VodPlayFrom string `json:"vod_play_from"` // 播放来源
}

// FilmClass 影视分类信息结构体
type FilmClass struct {
	TypeID   int64  `json:"type_id"`   // 分类ID
	TypePid  int64  `json:"type_pid"`  // 父级ID
	TypeName string `json:"type_name"` // 类型名称
}

// CommonPage 影视列表接口分页数据结构体
type CommonPage struct {
	Code      int    `json:"code"`      // 响应状态码
	Msg       string `json:"msg"`       // 数据类型
	Page      any    `json:"page"`      // 页码
	PageCount int    `json:"pagecount"` // 总页数
	Limit     any    `json:"limit"`     // 每页数据量
	Total     int    `json:"total"`     // 总数据量
}

// 保存到redis
func SaveFilmClass(list []FilmClass) error {
	data, _ := json.Marshal(list)
	return db.Rdb.Set(db.Cxt, config.FilmClassKey, data, config.ResourceExpired).Err()
}
