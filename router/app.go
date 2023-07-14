package router

//本文件是用于网关 路由相关的
import (
	"github.com/gin-gonic/gin"
	"main/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/index", service.GetIndex)

	r.GET("/user/getUserList", service.GetUserList)
	return r
}
