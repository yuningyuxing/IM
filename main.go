package main

import (
	"main/router"
	"main/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	r := router.Router()
	r.Run(":9090")
}
