package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//这里用来描述每个请求 并且让他可以在swagger上有所体现

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "welcome my home",
	})
}
