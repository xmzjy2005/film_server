package main

import (
	"film_server/config"
	"film_server/plugin/SystemInit"
	"film_server/plugin/db"
	"film_server/router"
	"fmt"
)

func init() {
	//初始化mysql
	err := db.InitMysql()
	if err != nil {
		panic(err)
	}
	//初始化redis
	err = db.InitRedisConn()
	if err != nil {
		panic(err)
	}
}
func main() {
	//这个没有初始化命令，直接写在这里
	DefaultDataInit()
	//开启路由监听
	r := router.SetupRouter()
	//开启监听
	r.Run(fmt.Sprintf(":%s", config.ListenerPort))
}
func DefaultDataInit() {
	//初始化影视资源来源
	SystemInit.SpiderInit()
}
