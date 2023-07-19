package service

import (
	"github.com/gin-gonic/gin"
	"html/template"
)

//这里用来描述每个请求 并且让他可以在swagger上有所体现

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")
	//c.JSON(http.StatusOK, gin.H{
	//	"message": "welcome my home",
	//})
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
}
