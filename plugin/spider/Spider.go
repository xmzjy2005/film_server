package spider

import (
	"errors"
	"film_server/model/system"
	"film_server/plugin/common/conver"
	"film_server/plugin/common/util"
	"fmt"
	"log"
	"net/url"
)

var spiderCore = &jsonCollect{}

// 执行采集
func HandleCollect(id string, h int) error {
	//查找采集站点信息
	s := system.FindCollectSourceById(id)
	fmt.Println(s, "Print filmSource")
	if s == nil {

		return errors.New("Cannot find collect source file")
	}
	if !s.State {
		return errors.New("The site state is false")
	}
	//如果是主站且状态启用，则获取分类树
	if s.Grade == system.MasterCollect && s.State {
		if !system.ExistsCategoryTree() {
			CollectCategory(s)
		}
	}
	//生成请求结构体
	r := util.RequestInfo{Uri: s.Uri, Params: url.Values{}}
	if h == 0 {
		return errors.New("collect time hour cannot be zero")
	}
	if h > 0 {
		r.Params.Set("h", fmt.Sprint(h))
	}
	//获取分页采集的页数
	pageCount, err := spiderCore.GetPageCount(r)
	if err != nil {
		pageCount, err = spiderCore.GetPageCount(r)
		if err != nil {
			return err
		}
	}
	fmt.Println(pageCount)
	//通过采集类型执行不同的采集方法
	switch s.CollectType {
	case system.CollectVideo: //采集 视频
		//if s.Interval > 500 {
		//	for i := 1; i <= pageCount; i++ {
		//		collectFilm(s, h, i)
		//	}
		//	time.Sleep(time.Duration(s.Interval) * time.Millisecond)
		//} else if pageCount <= config.MAXGoroutine*2 {
		//	for i := 1; i <= pageCount; i++ {
		//		collectFilm(s, h, i)
		//	}
		//} else {
		//	ConcurrenPageSiper(pageCount, s, h, collectFilm)
		//}

		//测试，先执行一个分页的数据采集
		collectFilm(s, h, 1)
		//视频采集数据完成后，同步到mysql
		if s.Grade == system.MasterCollect {
			if h > 0 {
				system.SyncSearchInfo(1)
			} else {
				system.SyncSearchInfo(0)
			}
			//开启图片同步
			if s.SyncPictures {
				system.SyncFilmPicture()
			}
			//执行完清除首页缓存
			ClearCache()
		}
		break
	default:
		fmt.Println("其他采集类型开发中...")
	}
	return nil
}

// 影视采集，单分页
func collectFilm(s *system.FilmSource, h, pg int) {
	//生成请求参数
	r := util.RequestInfo{Uri: s.Uri, Params: url.Values{}}
	r.Params.Set("pg", fmt.Sprint(pg))
	if h > 0 {
		r.Params.Set("h", fmt.Sprint(h))
	}
	//获取影片list
	list, err := spiderCore.GetFilmDetail(r)
	//抓取失败，存入record
	if err != nil || len(list) <= 0 {
		fr := system.FailureRecord{OriginId: s.Id, OriginName: s.Name, Uri: s.Uri, CollectType: s.CollectType, PageNumber: pg, Hour: h, Cause: fmt.Sprintln(err), Status: 1}
		system.SaveFailureRecord(fr)
	}
	switch s.Grade {
	case system.MasterCollect: //主站点
		err := system.SaveDetails(list)
		if err != nil {
			log.Println("Save Detail error :", err)
		}
		if s.SyncPictures {
			vPicList := conver.ConvertVirtualPicture(list)
			err := system.SaveVirtualPic(vPicList)
			if err != nil {
				log.Println("save virtual picture err:", err)
			}
		}
		break

	case system.SlaveCollect: //从站点
		if err = system.SaveSitePlayList(s.Id, list); err != nil {
			log.Println("save playlist err : ", err)
		}
		break
	}
}

// 影视分类
func CollectCategory(s *system.FilmSource) {
	//获取分类树行数据
	categoryTree, err := spiderCore.GetCategoryTree(util.RequestInfo{Uri: s.Uri, Params: url.Values{}})
	if err != nil {
		log.Println("GetCategoryTree errr :", err)
		return
	}
	//保存tree到redis
	err = system.SaveCategoryTree(categoryTree)
	if err != nil {
		log.Println("Save Category Tree error ", err)
	}
}
