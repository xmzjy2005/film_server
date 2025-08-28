package conver

import (
	"film_server/model/collect"
	"film_server/model/system"
	"strings"
)

func GenCategoryTree(list []collect.FilmClass) *system.CategoryTree {
	tree := &system.CategoryTree{Category: &system.Category{Id: 0, Pid: -1, Name: "分类信息", Show: true}}
	temp := make(map[int64]*system.CategoryTree)
	temp[tree.Id] = tree
	for _, c := range list {
		//判断当前节点是否在temp
		_, ok := temp[c.TypeID]
		if ok {
			temp[c.TypeID].Category = &system.Category{Id: c.TypeID, Pid: c.TypePid, Name: c.TypeName, Show: true}
		} else {
			temp[c.TypeID] = &system.CategoryTree{Category: &system.Category{Id: c.TypeID, Pid: c.TypePid, Name: c.TypeName, Show: true}}
		}
		//查找父节点
		pCategory, ok := temp[c.TypePid]
		if !ok {
			//不存在 先占个位置
			temp[c.TypePid] = &system.CategoryTree{Category: &system.Category{Id: c.TypePid}}
		}
		pCategory.Children = append(pCategory.Children, temp[c.TypeID])
	}
	return tree
}

// 批量处理影片详情
func ConvertFilmDetails(details []collect.FilmDetail) []system.MovieDetail {
	var dl []system.MovieDetail
	for _, d := range details {
		dl = append(dl, ConvertFilmDetail(d))
	}
	return dl
}

func ConvertFilmDetail(detail collect.FilmDetail) system.MovieDetail {
	md := system.MovieDetail{
		Id:       detail.VodID,
		Cid:      detail.TypeID,
		Pid:      detail.TypeID1,
		Name:     detail.VodName,
		Picture:  detail.VodPic,
		DownFrom: detail.VodDownFrom,
		MovieDescriptor: system.MovieDescriptor{
			SubTitle:    detail.VodSub,
			CName:       detail.TypeName,
			EnName:      detail.VodEn,
			Initial:     detail.VodLetter,
			ClassTag:    detail.VodClass,
			Actor:       detail.VodActor,
			Director:    detail.VodDirector,
			Writer:      detail.VodWriter,
			Blurb:       detail.VodBlurb,
			Remarks:     detail.VodRemarks,
			ReleaseDate: detail.VodPubDate,
			Area:        detail.VodArea,
			Language:    detail.VodLang,
			Year:        detail.VodYear,
			State:       detail.VodState,
			UpdateTime:  detail.VodTime,
			AddTime:     detail.VodTimeAdd,
			DbId:        detail.VodDouBanID,
			DbScore:     detail.VodDouBanScore,
			Hits:        detail.VodHits,
			Content:     detail.VodContent,
		},
	}
	md.PlayFrom = strings.Split(detail.VodPlayFrom, detail.VodPlayNote)
	md.PlayList = GenFilmPlayList(detail.VodPlayURL, detail.VodPlayNote)
	md.DownloadList = GenFilmPlayList(detail.VodDownURL, detail.VodPlayNote)
	return md
}
func GenFilmPlayList(playUrl, separator string) [][]system.MovieUrlInfo {
	var res [][]system.MovieUrlInfo
	if separator != "" {
		for _, l := range strings.Split(playUrl, separator) {
			if strings.Contains(l, ".m3u8") || strings.Contains(l, ".mp4") {
				res = append(res, ConverPlayUrl(l))
			}
		}
	} else {
		if strings.Contains(playUrl, ".m3u8") || strings.Contains(playUrl, ".mp4") {
			res = append(res, ConverPlayUrl(playUrl))
		}
	}
	return res
}

// ConvertPlayUrl 将单个playFrom的播放地址字符串处理成列表形式
/*
第01集$https://vod.360zyx.vip/20250724/lcqIqUXc/index.m3u8#第02集$https://vod.360zyx.vip/20250724/VDlQW9TS/index.m3u8#第03集$https://vod.360zyx.vip/20250724/NbegHrka/index.m3u8#第04集$https://vod.360zyx.vip/20250725/OwLfr9Cg/index.m3u8#第05集$https://vod.360zyx.vip/20250728/cvNxGWz4/index.m3u8#第06集$https://vod.360zyx.vip/20250728/LLesbqLn/index.m3u8
*/
func ConverPlayUrl(playUrl string) []system.MovieUrlInfo {
	var l []system.MovieUrlInfo
	for _, p := range strings.Split(playUrl, "#") {
		if strings.Contains(p, "$") {
			l = append(l, system.MovieUrlInfo{
				Episode: strings.Split(p, "$")[0],
				Link:    strings.Split(p, "$")[1],
			})
		} else {
			l = append(l, system.MovieUrlInfo{
				Episode: "null",
				Link:    p,
			})
		}
	}
	return l
}

// 将影片详情转换为虚拟图片
func ConvertVirtualPicture(details []system.MovieDetail) []system.VirtualPicture {
	var l []system.VirtualPicture
	for _, d := range details {
		if len(d.Picture) > 0 {
			l = append(l, system.VirtualPicture{
				Id:   d.Id,
				Link: d.Picture,
			})
		}

	}
	return l
}
