package system

import (
	"encoding/json"
	"film_server/config"
	"film_server/plugin/common/util"
	"film_server/plugin/db"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// FileInfo 图片信息对象
type FileInfo struct {
	gorm.Model
	Link        string `json:"link"`        // 图片链接
	Uid         int    `json:"uid"`         // 上传人ID
	RelevanceId int64  `json:"relevanceId"` // 关联资源ID
	Type        int    `json:"type"`        // 文件类型 (0 影片封面, 1 用户头像)
	Fid         string `json:"fid"`         // 图片唯一标识, 通常为文件名
	FileType    string `json:"fileType"`    // 文件类型, txt, png, jpg
	//Size        int    `json:"size"`        // 文件大小
}

// 采集入站 到redis的存储对象
type VirtualPicture struct {
	Id   int64  `json:"id"`
	Link string `json:"link"`
}

// 保存照片信息到redis
func SaveVirtualPic(pl []VirtualPicture) error {
	var zl []redis.Z
	for _, p := range pl {
		m, _ := json.Marshal(p)
		zl = append(zl, redis.Z{
			Score:  float64(p.Id),
			Member: m,
		})
	}
	return db.Rdb.ZAdd(db.Cxt, config.VirtualPictureKey, zl...).Err()
}

// 将redis的信息保存到本地
func SyncFilmPicture() {
	//获取缓存中图片数量
	count := db.Rdb.ZCard(db.Cxt, config.VirtualPictureKey).Val()
	if count <= 0 {
		return
	}
	//扫描固定条数
	sl := db.Rdb.ZPopMax(db.Cxt, config.VirtualPictureKey, config.MaxScanCount).Val()
	if len(sl) <= 0 {
		return
	}

	for _, s := range sl {
		vp := VirtualPicture{}
		_ = json.Unmarshal([]byte(s.Member.(string)), &vp)
		//判断当前影片是否同步过图片
		if ExistFileInfoByRid(vp.Id) {
			continue
		}
		//将图片保存到服务器中
		util.SaveOnlineFile(vp.Link, config.FilmPictureUploadDir)
	}
}
func ExistFileInfoByRid(rid int64) bool {
	var count int64
	db.Mdb.Model(&FileInfo{}).Where("relevance_id = ?", rid).Count(&count)
	return count > 0
}
