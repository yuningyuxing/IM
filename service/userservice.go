package service

import (
	"github.com/gin-gonic/gin"
	"main/models"
	"net/http"
)

//注意这不是注释  '//@'是用来添加特定地注释 这种注释通常被称为Swagger注解 它们用于为API定义添加元数据和描述
//swagger可以解析这些注解 并生成相应地API文档

// GetUserList
// @Tags 首页
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}
