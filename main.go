package main

import (
	"github.com/spf13/viper"
	"main/models"
	"main/router"
	"main/utils"
	"time"
)

// 初始化定时器  定时器的作用是定时清理超时连接
func InitTimer() {
	//给写好的定时器传入参数
	utils.Timer(
		time.Duration(viper.GetInt("timeout.DelayHeartbeat"))*time.Second,
		time.Duration(viper.GetInt("timeout.HeartbeatHz"))*time.Second,
		models.CleanConnection,
		"")
}

func main() {

	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	InitTimer()
	r := router.Router()
	r.Run(viper.GetString("port.server"))
}
