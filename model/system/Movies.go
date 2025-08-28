package system

import (
	"encoding/json"
	"film_server/config"
	"film_server/plugin/db"
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MovieDetail 影片详情信息
type MovieDetail struct {
	Id       int64    `json:"id"`       //影片Id
	Cid      int64    `json:"cid"`      //分类ID
	Pid      int64    `json:"pid"`      //一级分类ID
	Name     string   `json:"name"`     //片名
	Picture  string   `json:"picture"`  //简介图片
	PlayFrom []string `json:"playFrom"` // 播放来源
	DownFrom string   `json:"DownFrom"` //下载来源 例: http
	//PlaySeparator   string              `json:"playSeparator"` // 播放信息分隔符
	PlayList        [][]MovieUrlInfo    `json:"playList"`     //播放地址url
	DownloadList    [][]MovieUrlInfo    `json:"downloadList"` // 下载url地址
	MovieDescriptor `json:"descriptor"` //影片描述信息
}

// MovieUrlInfo 影视资源url信息
type MovieUrlInfo struct {
	Episode string `json:"episode"` // 集数
	Link    string `json:"link"`    // 播放地址
}

// MovieDescriptor 影片详情介绍信息
type MovieDescriptor struct {
	SubTitle    string `json:"subTitle"`    //子标题
	CName       string `json:"cName"`       //分类名称
	EnName      string `json:"enName"`      //英文名
	Initial     string `json:"initial"`     //首字母
	ClassTag    string `json:"classTag"`    //分类标签
	Actor       string `json:"actor"`       //主演
	Director    string `json:"director"`    //导演
	Writer      string `json:"writer"`      //作者
	Blurb       string `json:"blurb"`       //简介, 残缺,不建议使用
	Remarks     string `json:"remarks"`     // 更新情况
	ReleaseDate string `json:"releaseDate"` //上映时间
	Area        string `json:"area"`        // 地区
	Language    string `json:"language"`    //语言
	Year        string `json:"year"`        //年份
	State       string `json:"state"`       //影片状态 正片|预告...
	UpdateTime  string `json:"updateTime"`  //更新时间
	AddTime     int64  `json:"addTime"`     //资源添加时间戳
	DbId        int64  `json:"dbId"`        //豆瓣id
	DbScore     string `json:"dbScore"`     // 豆瓣评分
	Hits        int64  `json:"hits"`        //影片热度
	Content     string `json:"content"`     //内容简介
}

// MovieBasicInfo 影片基本信息
type MovieBasicInfo struct {
	Id       int64  `json:"id"`       //影片Id
	Cid      int64  `json:"cid"`      //分类ID
	Pid      int64  `json:"pid"`      //一级分类ID
	Name     string `json:"name"`     //片名
	SubTitle string `json:"subTitle"` //子标题
	CName    string `json:"cName"`    //分类名称
	State    string `json:"state"`    //影片状态 正片|预告...
	Picture  string `json:"picture"`  //简介图片
	Actor    string `json:"actor"`    //主演
	Director string `json:"director"` //导演
	Blurb    string `json:"blurb"`    //简介, 不完整
	Remarks  string `json:"remarks"`  // 更新情况
	Area     string `json:"area"`     // 地区
	Year     string `json:"year"`     //年份
}

func SaveDetails(list []MovieDetail) error {
	var err error
	var searchList []SearchInfo
	for _, detail := range list {
		data, _ := json.Marshal(detail)
		//存入redis
		err = db.Rdb.Set(db.Cxt, fmt.Sprintf(config.MovieDetailKey, detail.Cid, detail.Id), data, config.FilmExpired).Err()
		//保存基本信息到redis
		SaveMovieBasicInfo(detail)
		//转换detail信息为searchInfo
		searchInfo := ConvertSearchInfo(detail)
		searchList = append(searchList, searchInfo)
		//存储检索信息
		SaveSearchTag(searchInfo)
	}
	//保存一份search信息到Mysql
	RdbSaveSearchInfo(searchList)
	return err
}
func SaveMovieBasicInfo(detail MovieDetail) {
	basicInfo := MovieBasicInfo{
		Id:       detail.Id,
		Cid:      detail.Cid,
		Pid:      detail.Pid,
		Name:     detail.Name,
		SubTitle: detail.SubTitle,
		CName:    detail.CName,
		State:    detail.State,
		Picture:  detail.Picture,
		Actor:    detail.Actor,
		Director: detail.Director,
		Blurb:    detail.Blurb,
		Remarks:  detail.Remarks,
		Area:     detail.Area,
		Year:     detail.Year,
	}
	data, _ := json.Marshal(basicInfo)
	_ = db.Rdb.Set(db.Cxt, fmt.Sprintf(config.MovieBasicInfoKey, detail.Cid, detail.Id), data, config.FilmExpired).Err()
}
func ConvertSearchInfo(detail MovieDetail) SearchInfo {
	year, err := strconv.ParseInt(regexp.MustCompile("[1-9][0-9]{3}").FindString(detail.ReleaseDate), 10, 64)
	if err != nil {
		year = 0
	}
	score, _ := strconv.ParseFloat(detail.DbScore, 64)
	stamp, _ := time.ParseInLocation(time.DateTime, detail.UpdateTime, time.Local)
	return SearchInfo{
		Mid:          detail.Id,
		Cid:          detail.Cid,
		Pid:          detail.Pid,
		Name:         detail.Name,
		SubTitle:     detail.SubTitle,
		CName:        detail.CName,
		ClassTag:     detail.ClassTag,
		Area:         detail.Area,
		Language:     detail.Language,
		Year:         year,
		Initial:      detail.Initial,
		Score:        score,
		Hits:         detail.Hits,
		UpdateStamp:  stamp.Unix(),
		State:        detail.State,
		Remarks:      detail.Remarks,
		ReleaseStamp: detail.AddTime,
	}
}

// 保存从站的返回信息
func SaveSitePlayList(id string, list []MovieDetail) (err error) {
	if len(list) <= 0 {
		return nil
	}
	res := make(map[string]string)
	for _, d := range list {
		fmt.Println(d)
		//集数跟连接
		if len(d.PlayList) > 0 {
			//第一集压缩？
			data, _ := json.Marshal(d.PlayList[0])
			if strings.Contains(d.CName, "解说") {
				continue
			}
			//如果豆瓣id有值，则用豆瓣id再存储一次
			if d.DbId != 0 {
				res[GenerateHashKey(d.DbId)] = string(data)
			}
			res[GenerateHashKey(d.Name)] = string(data)
		}
	}
	if len(res) > 0 {
		err = db.Rdb.HMSet(db.Cxt, fmt.Sprintf(config.MultipleSiteDetail, id), res).Err()
	}

	return
}

/*
	对附属播放源入库时的name|dbID进行处理,保证唯一性
1. 去除name中的所有空格
2. 去除name中含有的别名～.*～
3. 去除name首尾的标点符号
4. 将处理完成后的name转化为hash值作为存储时的key
*/
// GenerateHashKey 存储播放源信息时对影片名称进行处理, 提高各站点间同一影片的匹配度
func GenerateHashKey[k string | ~int | int64](key k) string {
	mName := fmt.Sprint(key)
	//去除所有空格
	mName = regexp.MustCompile(`\s`).ReplaceAllString(mName, "")
	//去除结尾~.*~
	mName = regexp.MustCompile(`~.*~$`).ReplaceAllString(mName, "")
	//去除首位的标点符号
	mName = regexp.MustCompile(`^[[:punct:]]+ | [[:punct:]]+$`).ReplaceAllString(mName, "")
	//去除季后面的
	mName = regexp.MustCompile(`季.*`).ReplaceAllString(mName, "季")
	//hash计算
	h := fnv.New32a()
	_, err := h.Write([]byte(mName))
	if err != nil {
		return ""
	}
	return fmt.Sprint(h.Sum32())
}
