package system

import (
	"encoding/json"
	"film_server/config"
	"film_server/plugin/db"
	"github.com/redis/go-redis/v9"
	"log"
)

type CollectResultModel int
type SourceGrade int
type ResourceType int

// FilmSource 影视站点信息保存结构体
type FilmSource struct {
	Id           string             `json:"id"`           // 唯一ID
	Name         string             `json:"name"`         // 采集站点备注名
	Uri          string             `json:"uri"`          // 采集链接
	ResultModel  CollectResultModel `json:"resultModel"`  // 接口返回类型, json || xml
	Grade        SourceGrade        `json:"grade"`        // 采集站等级 主站点 || 附属站
	SyncPictures bool               `json:"syncPictures"` // 是否同步图片到服务器
	CollectType  ResourceType       `json:"collectType"`  // 采集资源类型
	State        bool               `json:"state"`        // 是否启用
	Interval     int                `json:"interval"`     // 采集时间间隔 单位/ms
}

const (
	CollectVideo = iota
	CollectArticle
	CollectActor
	CollectRole
	CollectWebSite
)
const (
	MasterCollect SourceGrade = iota
	SlaveCollect
)
const (
	JsonResult CollectResultModel = iota
	XmlResult
)

func GetCollectSourceList() []FilmSource {
	l, err := db.Rdb.ZRange(db.Cxt, config.FilmSourceListKey, 0, -1).Result()
	if err != nil {
		log.Println(err)
		return nil
	}
	return getCollectSource(l)

}

// 格式化redis取值
func getCollectSource(sl []string) []FilmSource {
	var l []FilmSource
	for _, s := range sl {
		f := FilmSource{}
		_ = json.Unmarshal([]byte(s), &f)
		l = append(l, f)
	}
	return l
}

// 判断是否redis有数据
func ExistCollectSourceList() bool {
	if db.Rdb.Exists(db.Cxt, config.FilmSourceListKey).Val() == 0 {
		return false
	}
	return true
}

// 存储采集列表
func SaveCollectSourceList(list []FilmSource) error {
	var zl []redis.Z
	for _, v := range list {
		s, _ := json.Marshal(v)
		zl = append(zl, redis.Z{
			Score:  float64(v.Grade),
			Member: s,
		})
	}
	return db.Rdb.ZAdd(db.Cxt, config.FilmSourceListKey, zl...).Err()
}

// 通过id查找资源站信息
func FindCollectSourceById(id string) *FilmSource {
	for _, v := range GetCollectSourceList() {
		if v.Id == id {
			return &v
		}
	}
	return nil
}
