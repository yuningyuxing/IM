package router

//本文件是用于网关 路由相关的
import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"main/docs"
	"main/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//加载静态资源
	//设置静态文件目录映射
	r.Static("/asset", "asset/")
	//加载HTML模板文件
	// 这样，在使用模板引擎渲染页面时，可以直接通过模板文件的相对路径引用模板，而不需要指定完整的文件路径。
	// 例如，使用{{ template "subdir/template.html" . }}可以引用"views/subdir/template.html"文件作为模板
	r.LoadHTMLGlob("views/**/*")
	//首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/chat", service.Chat)
	r.GET("/toChat", service.ToChat)
	//加载好友
	r.POST("/searchFriends", service.SearchFriends)
	//加载群列表
	r.POST("/contact/loadcommunity", service.LoadCommunity)

	//用户模块
	r.POST("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)

	//发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	//上传文件
	r.POST("/attach/upload", service.Upload)
	//添加好友
	r.POST("/contact/addfriend", service.AddFriend)
	//创建群
	r.POST("/contact/createCommunity", service.CreateCommunity)
	//加群
	r.POST("/contact/joinGroup", service.JoinGroup)

	return r
}
