package main

import (
	"main/router"
	"main/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	r := router.Router()
	r.Run(":9090")
}
