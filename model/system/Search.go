package system

import (
	"encoding/json"
	"film_server/config"
	"film_server/plugin/db"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"regexp"
	"strings"
	"time"
)

type SearchInfo struct {
	gorm.Model
	Mid          int64   `json:"mid"`          //影片ID gorm:"uniqueIndex:idx_mid"
	Cid          int64   `json:"cid"`          //分类ID
	Pid          int64   `json:"pid"`          //上级分类ID
	Name         string  `json:"name"`         // 片名
	SubTitle     string  `json:"subTitle"`     // 影片子标题
	CName        string  `json:"cName"`        // 分类名称
	ClassTag     string  `json:"classTag"`     //类型标签
	Area         string  `json:"area"`         // 地区
	Language     string  `json:"language"`     // 语言
	Year         int64   `json:"year"`         // 年份
	Initial      string  `json:"initial"`      // 首字母
	Score        float64 `json:"score"`        //评分
	UpdateStamp  int64   `json:"updateStamp"`  // 更新时间
	Hits         int64   `json:"hits"`         // 热度排行
	State        string  `json:"state"`        //状态 正片|预告
	Remarks      string  `json:"remarks"`      // 完结 | 更新至x集
	ReleaseStamp int64   `json:"releaseStamp"` //上映时间 时间戳
}

func (s *SearchInfo) TableName() string {
	return config.SearchTableName
}

// 保存检索信息
func SaveSearchTag(search SearchInfo) {
	// 声明用于存储采集的影片的分类检索信息
	// Redis中的记录形式 Search:SearchKeys:Pid1:Title Hash
	// Redis中的记录形式 Search:SearchKeys:Pid1:xxx Hash
	key := fmt.Sprintf(config.SearchTitle, search.Pid)
	searchMap := db.Rdb.HGetAll(db.Cxt, key).Val()
	if len(searchMap) == 0 {
		searchMap = make(map[string]string)
		searchMap["Category"] = "类型"
		searchMap["Plot"] = "剧情"
		searchMap["Area"] = "地区"
		searchMap["Language"] = "语言"
		searchMap["Year"] = "年份"
		searchMap["Initial"] = "首字母"
		searchMap["Sort"] = "排序"
		db.Rdb.HMSet(db.Cxt, key, searchMap)
	}
	for k, _ := range searchMap {
		tagKey := fmt.Sprintf(config.SearchTag, search.Pid, k)
		tagCount := db.Rdb.ZCard(db.Cxt, tagKey).Val()
		switch k {
		case "Category":
			if tagCount == 0 {
				for _, t := range GetChildrenTree(search.Pid) {
					db.Rdb.ZAdd(db.Cxt, tagKey, redis.Z{
						Score:  float64(-t.Id),
						Member: fmt.Sprintf("%v:%v", t.Name, t.Id),
					})
				}
			}
		case "Year":
			if tagCount == 0 {
				curYear := time.Now().Year()
				for i := 0; i < 12; i++ {
					db.Rdb.ZAdd(db.Cxt, tagKey, redis.Z{
						Score:  float64(curYear - i),
						Member: fmt.Sprintf("%v:%v", curYear-i, curYear-i),
					})
				}
			}
		case "Initial":
			if tagCount == 0 {
				for i := 65; i <= 90; i++ {
					db.Rdb.ZAdd(db.Cxt, tagKey, redis.Z{
						Score:  float64(90 - i),
						Member: fmt.Sprintf("%v:%v", i, i),
					})
				}
			}
		case "Sort":
			if tagCount == 0 {
				tags := []redis.Z{
					{Score: 3, Member: "时间排序:update_stamp"},
					{Score: 2, Member: "人气排序:hits"},
					{Score: 1, Member: "评分排序:score"},
					{Score: 0, Member: "最新上映:release_stamp"},
				}
				db.Rdb.ZAdd(db.Cxt, tagKey, tags...)
			}
		case "Plot": //类型标签
			HandleSearchTag(search.ClassTag, tagKey)
		case "Area":
			HandleSearchTag(search.ClassTag, tagKey)
		case "Language":
			HandleSearchTag(search.Language, tagKey)
		default:
			break
		}

	}
}

/*
@name 存入标签redis
@param preTags string 例如 动作/喜剧
@param k string Search:Pid%d:%s 例如Search:0:Plot
*/
func HandleSearchTag(preTags string, k string) {
	//剔除不要的
	preTags = regexp.MustCompile(`[\s\n\r]+`).ReplaceAllString(preTags, "")
	//切割存入redis分数
	f := func(sep string) {
		for _, t := range strings.Split(preTags, sep) {
			//获取现有分数
			score := db.Rdb.ZScore(db.Cxt, k, fmt.Sprintf("%v:%v", t, t)).Val()
			//现有分数加1
			db.Rdb.ZAdd(db.Cxt, k, redis.Z{
				Score:  score + 1,
				Member: fmt.Sprintf("%v:%v", t, t),
			})
		}
	}
	if strings.Contains(preTags, "/") {
		f("/")
	} else if strings.Contains(preTags, ",") {
		f(",")
	} else if strings.Contains(preTags, "，") {
		f("，")
	} else if strings.Contains(preTags, "、") {
		f("、")
	} else {
		if len(preTags) == 0 {

		} else if preTags == "其他" {
			db.Rdb.ZAdd(db.Cxt, k, redis.Z{
				Score:  0,
				Member: fmt.Sprintf("%v:%v", preTags, preTags),
			})
		} else {
			score := db.Rdb.ZScore(db.Cxt, k, fmt.Sprintf("%v:%v", preTags, preTags)).Val()
			db.Rdb.ZAdd(db.Cxt, k, redis.Z{
				Score:  score + 1,
				Member: fmt.Sprintf("%v:%v", preTags, preTags),
			})
		}
	}
}

// 存储搜索信息到redis
func RdbSaveSearchInfo(list []SearchInfo) {
	var members []redis.Z
	for _, info := range list {
		data, _ := json.Marshal(info)
		members = append(members, redis.Z{
			Score:  float64(info.Mid),
			Member: data,
		})
	}
	db.Rdb.ZAdd(db.Cxt, config.SearchInfoTemp, members...)
}

// 将影片存储到mysql
func SyncSearchInfo(model int) {
	if model == 0 { //传入的time ==-1的才有，即采集方式为采集全部
		//重置search表
		ResetSearchTable()
		//批量添加 searchInfo
		SearchInfoMdb(model)
		//给表添加索引
		AddSearchIndex()
	} else if model == 1 {
		//批量更新或者添加
		SearchInfoMdb(model)
	}
}
func AddSearchIndex() {
	var s SearchInfo
	tableName := s.TableName()
	db.Mdb.Exec(fmt.Sprintf("CREATE UNIQUE INDEX idx_mid on %s (mid)", tableName))
	db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_time ON %s (update_stamp DESC)", tableName))
	db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_hits ON %s (hits DESC)", tableName))
	db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_score ON %s (score DESC)", tableName))
	db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_release ON %s (release_stamp DESC)", tableName))
	db.Mdb.Exec(fmt.Sprintf("CREATE INDEX idx_year ON %s (year DESC)", tableName))
}

// 扫描redis 并且存入mysql当中
func SearchInfoMdb(model int) {
	count := db.Rdb.ZCard(db.Cxt, config.SearchInfoTemp).Val()
	if count <= 0 {
		return
	}
	list := db.Rdb.ZPopMax(db.Cxt, config.SearchInfoTemp, config.MaxScanCount).Val()
	if len(list) <= 0 {
		return
	}
	var sl []SearchInfo
	for _, s := range list {
		info := SearchInfo{}
		_ = json.Unmarshal([]byte(s.Member.(string)), &info)
		sl = append(sl, info)
	}
	if model == 0 {
		//全部更新
		BatchSave(sl)
	} else {
		BatchSaveOrUpdate(sl)
	}
}

// 插入或者更新
func BatchSaveOrUpdate(list []SearchInfo) {
	tx := db.Mdb.Begin()
	for _, info := range list {
		var count int64
		//查下是否存在记录
		tx.Model(&SearchInfo{}).Where("mid", info.Mid).Count(&count)
		if count > 0 {
			//存在记录更新
			err := tx.Model(&SearchInfo{}).Where("mid", info.Mid).Updates(SearchInfo{
				UpdateStamp: info.UpdateStamp, Hits: info.Hits, State: info.State,
				Remarks: info.Remarks, Score: info.Score, ReleaseStamp: info.ReleaseStamp,
			}).Error
			if err != nil {
				tx.Rollback()
			}
		} else {
			err := tx.Create(&info).Error
			if err != nil {
				tx.Rollback()
			}
			//插入成功后保存一份检索tag
			BatchHandleSearchTag(info)
		}
	}
	tx.Commit()
}
func BatchSave(sl []SearchInfo) {
	tx := db.Mdb.Begin()
	//如果程序panic 则用recover捕获，并且回滚数据
	defer func() {
		r := recover()
		if r != nil {
			tx.Rollback()
		}
	}()
	err := tx.CreateInBatches(sl, len(sl)).Error
	if err != nil {
		tx.Rollback()
	}
	BatchHandleSearchTag(sl...)
	tx.Commit()
}
func BatchHandleSearchTag(infos ...SearchInfo) {
	for _, info := range infos {
		SaveSearchTag(info)
	}
}
func ResetSearchTable() {
	var s SearchInfo
	db.Mdb.Exec(fmt.Sprintf("drop table if exists %s", s.TableName()))
	CreateSearchTable()
}
func CreateSearchTable() {
	if !ExistSearchTable() {
		err := db.Mdb.AutoMigrate(&SearchInfo{})
		if err != nil {
			log.Println("Create table searchInfoFail:", err)
		}
	}
}
func ExistSearchTable() bool {
	return db.Mdb.Migrator().HasTable(&SearchInfo{})
}
