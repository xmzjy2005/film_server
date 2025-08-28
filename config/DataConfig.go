package config

import "time"

// ------System config ------
const (
	ListenerPort         = "3601"
	MAXGoroutine         = 10
	FilmPictureUploadDir = "./static/upload/gallery" //保存影片的封面图片的本地目录 这是哪里？？
)
const (
	// SearchTableName 存放检索信息的数据表名
	SearchTableName        = "search"
	UserTableName          = "users"
	UserIdInitialVal       = 10000
	FileTableName          = "files"
	FailureRecordTableName = "failure_records"

	MysqlDsn = "root:123456@(192.168.108.21:3306)/FilmSite?charset=utf8mb4&parseTime=True&loc=Local"

	RedisAddr     = `192.168.108.21:6379`
	RedisPassword = `root`
	RedisDBNo     = 4
)
const (
	AuthUserClaims = "UserClaims"
)

// --------后台管理key---------------
const (
	// FilmSourceListKey 采集 API 信息列表key
	FilmSourceListKey = "Config:Collect:FilmSource"
)

// -------------------------redis key-----------------------------------
const (
	CategoryTreeKey   = "CategoryTree"
	FilmExpired       = time.Hour * 24 * 365 * 10
	MovieDetailKey    = "MovieDetail:Cid%d:Id:%d"
	MovieBasicInfoKey = "MovieBasicInfo:Cid%d:Id:%d"
	//搜索redis 存Pid 上级目录id
	SearchTitle    = "Search:Pid%d:Title"
	SearchTag      = "Search:Pid%d:%s"
	SearchInfoTemp = "Search:SearchInfoTemp"
	//待同步的图片临时存储
	VirtualPictureKey = "VirtualPicture"
	//多站点（从站点）存储影片信息
	MultipleSiteDetail = "MultipleSource:%s"
	//扫描redis到mysql一次性最大数量
	MaxScanCount = 300
)
