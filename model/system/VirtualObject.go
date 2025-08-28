package system

// UserInfoVo 用户信息返回对象
type UserInfoVo struct {
	Id       uint   `json:"id"`
	UserName string `json:"userName"` // 用户名
	Email    string `json:"email"`    // 邮箱
	Gender   int    `json:"gender"`   // 性别
	NickName string `json:"nickName"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
	Status   int    `json:"status"`   // 状态
}

// CollectParams 数据采集所需要的参数
type CollectParams struct {
	Id    string   `json:"id"`    // 资源站id
	Ids   []string `json:"ids"`   // 资源站id列表
	Time  int      `json:"time"`  // 采集时长
	Batch bool     `json:"batch"` // 是否批量执行
}
