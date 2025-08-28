package SystemInit

import (
	"film_server/model/system"
	"film_server/plugin/util"
	"log"
)

func SpiderInit() {
	FilmSourceInit()
}
func FilmSourceInit() {
	if system.ExistCollectSourceList() == true {
		return
	}
	l := []system.FilmSource{
		{Id: util.GenerateSalt(), Name: "HD(LZ)", Uri: `https://cj.lziapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		{Id: util.GenerateSalt(), Name: "HD(BF)", Uri: `https://bfzyapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false, Interval: 2500},
		{Id: util.GenerateSalt(), Name: "HD(FF)", Uri: `http://cj.ffzyapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		{Id: util.GenerateSalt(), Name: "HD(OK)", Uri: `https://okzyapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		{Id: util.GenerateSalt(), Name: "HD(HM)", Uri: `https://json.heimuer.xyz/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		{Id: util.GenerateSalt(), Name: "HD(LY)", Uri: `https://360zy.com/api.php/provide/vod/at/json`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		{Id: util.GenerateSalt(), Name: "HD(SN)", Uri: `https://suoniapi.com/api.php/provide/vod/from/snm3u8/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false, Interval: 2000},
		{Id: util.GenerateSalt(), Name: "HD(DB)", Uri: `https://caiji.dbzy.tv/api.php/provide/vod/from/dbm3u8/at/josn/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		{Id: util.GenerateSalt(), Name: "HD(IK)", Uri: `https://ikunzyapi.com/api.php/provide/vod/at/json`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		//{Id: util.GenerateSalt(), Name: "WX(T2)", Uri: `https://api.wuxianzy.net/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		//{Id: util.GenerateSalt(), Name: "OK(BK)", Uri: `https://api.okzy.org/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		//{Id: util.GenerateSalt(), Name: "HD(HW)", Uri: `https://cjhwba.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false, Interval: 3000},
		//{Id: util.GenerateSalt(), Name: "HD(lzBk)", Uri: `https://cj.lzcaiji.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		//{Id: util.GenerateSalt(), Name: "HD(fs)", Uri: `https://www.feisuzyapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
		//{Id: util.GenerateSalt(), Name: "HD(bfBk)", Uri: `http://app.bfzyapi.com/api.php/provide/vod/`, ResultModel: system.JsonResult, Grade: system.SlaveCollect, SyncPictures: false, CollectType: system.CollectVideo, State: false},
	}
	err := system.SaveCollectSourceList(l)
	if err != nil {
		log.Println("save filmList to redis error", err)
	}
}
