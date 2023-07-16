package service

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"main/models"
	"main/utils"
	"math/rand"
	"net/http"
	"strconv"
)

//注意这不是注释  '//@'是用来添加特定地注释 这种注释通常被称为Swagger注解 它们用于为API定义添加元数据和描述
//swagger可以解析这些注解 并生成相应地API文档

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0表示成功 -1表示失败
		"message": "获取全部用户成功",
		"data":    data,
	})
}

//summary = 总结，概要     param = 参数

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	//创建新用户
	//先创建一个用户模板 然后通过c获取参数
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")

	//设置盐值 随机数
	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1, //0表示成功 -1表示失败
			"message": "用户名已被使用",
			"data":    data,
		})
		return
	}
	//判断密码
	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1, //0表示成功 -1表示失败
			"message": "两次密码不一致",
			"data":    data,
		})
		return
	}
	user.Salt = salt
	user.PassWord = utils.MakePassword(password, salt)
	//去数据库创建用户
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0表示成功 -1表示失败
		"message": "用户新增成功",
		"data":    data,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	//删除用户  根据ID
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0表示成功 -1表示失败
		"message": "用户删除成功",
		"data":    user,
	})
}

// UpdateUser
// @Summary 更新用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	//更新  根据ID
	user := models.UserBasic{}
	//注意PostForm得到的是string 需要转化成int 然后再转化成uint
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	//调用govalidator库 来校验参数是否与规则匹配 规则在定义哪里时制定
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1, //0表示成功 -1表示失败
			"message": "修改参数不匹配",
			"data":    user,
		})
		return
	}

	models.UpdateUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0表示成功 -1表示失败
		"message": "用户更新成功",
		"data":    user,
	})
}

// FindUserByNameAndPwd
// @Summary 用户登录
// @Tags 用户模块
// @param name query string false "name"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	user := models.UserBasic{}
	name := c.Query("name")
	password := c.Query("password")
	user = models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //0表示成功 -1表示失败
			"message": "用户不存在",
			"data":    user,
		})
		return
	}
	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, //0表示成功 -1表示失败
			"message": "密码错误",
			"data":    user,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data := models.FindUserByNameAndPwd(name, pwd)
	c.JSON(200, gin.H{
		"code":    0, //0表示成功 -1表示失败
		"message": "用户登录成功",
		"data":    data,
	})
}
