package main

import (
	"github.com/spf13/viper"
	"main/models"
	"main/router"
	"main/utils"
	"time"
)

// 初始化定时器
func InitTimer() {
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
