package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"main/utils"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Upload(c *gin.Context) {
	UploadLocal(c)
}

// 上传文件到本地
func UploadLocal(c *gin.Context) {
	//这里获取ResponseWriter和Request对象 等会用
	w := c.Writer
	req := c.Request
	//从表单中获取文件对象和文件头 以及错误
	//head包含文件名 大小 内容类型等
	srcFile, head, err := req.FormFile("file")
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	//设置一个默认后缀
	suffix := ".png"
	//获取文件名
	ofilName := head.Filename
	//切分文件名 获取文件后缀名
	tem := strings.Split(ofilName, ".")
	if len(tem) > 1 {
		//获取文件后缀名
		suffix = "." + tem[len(tem)-1]
	}
	//构建新的文件名 包含时间戳和随机数  避免文件名冲突

	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	println("ssss")
	println(fileName)
	//创建目标文件并打开  此时文件就是dstFile 当然这时候还没东西
	dstFile, err := os.Create("./asset/upload/" + fileName)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	//将源文件内容拷贝到目标文件中
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	//构建URL 返回成功响应
	url := "./asset/upload/" + fileName
	utils.RespOK(w, url, "发送图片成功")
}
