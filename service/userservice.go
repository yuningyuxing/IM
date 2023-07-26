package service

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"main/models"
	"main/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"
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
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param repassword formData string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	//创建新用户
	//先创建一个用户模板 然后通过c获取参数
	user := models.UserBasic{}
	//POST用FormValue
	//user.Name = c.Query("name")
	//password := c.Query("password")
	//repassword := c.Query("repassword")
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("Identity")
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

	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(-1, gin.H{
			"code":    -1, //0表示成功 -1表示失败
			"message": "用户名或密码不能为空",
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
	utils.RespOK(c.Writer, data, "创建成功")
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id formData string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [POST]
func DeleteUser(c *gin.Context) {
	//删除用户  根据ID
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	fmt.Println(id)
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
	//name := c.Query("name")
	//password := c.Query("password")
	//注意Request.FormValue可以获取URL和POST表单中的参数 而Query只能获取到URL中的
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
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

// 防止跨域站点伪造请求
// upGrader是一个websocket.Upgrader类型的变量 用于将HTTP连接升级为WebSocket连接
var upGrader = websocket.Upgrader{
	//在websocket连接过程中 客户端会发送一个Origin头部字段 用于表示请求来源
	//checkorigin函数会被调用来检查该来源是否合法
	//默认情况下checkorigin的函数是nil 既不进行来源检查 允许来自任意来源的连接 为了增加安全性可以自定义checkorigin函数
	//这里被始终定义为返回true 这代表不进行源检查
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 用于处理websocket连接和发送消息
func SendMsg(c *gin.Context) {
	//将HTTP连接升级为websocket连接 并获取升级后的websocket连接
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//延迟执行关闭websocket连接的操作
	//(ws)是传入参数 因为匿名函数无法直接访问外部函数的变量
	//只有在匿名函数后面加上括号(可以传参)这样他才会执行
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	//调用MsgHandler函数处理接受到的消息
	MsgHandler(c, ws)
}

// 加载用户的缓存
func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := models.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOKList(c.Writer, "ok", res)
}

// 用于处理websocket连接的消息
func MsgHandler(c *gin.Context, ws *websocket.Conn) {
	//循环 接受redis订阅频道的消息并处理
	for {
		//从订阅的频道中接受消息
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println("MsgHandler 发送失败 ", err)
		}
		//获取当前时间并格式化
		tm := time.Now().Format("2006-01-02 15:04:05")
		//构造要发送的消息内容
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		//fmt.Println(m)
		//向websocket连接发送消息 消息类型为1
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

// 查询好友列表
func SearchFriends(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriends(uint(userId))
	utils.RespOKList(c.Writer, users, len(users))
}

// 添加好友
func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetName := c.Request.FormValue("targetName")
	code, msg := models.AddFriend(uint(userId), targetName)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 创建群
func CreateCommunity(c *gin.Context) {
	community := models.Community{}
	community.Name = c.Request.PostFormValue("name")
	ownerId, _ := strconv.Atoi(c.Request.PostFormValue("ownerId"))
	community.OwnerId = uint(ownerId)
	code, msg := models.CreatCommunity(community)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 加载群列表
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.PostFormValue("ownerId"))
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespOKList(c.Writer, data, len(data))
	} else {
		utils.RespFail(c.Writer, msg)
	}

}

// 加群
func JoinGroup(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	groupName := c.Request.FormValue("comId")
	code, msg := models.JoinGroup(uint(userId), groupName)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 根据ID找到用户
func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	data := models.FindByID(uint(userId))
	utils.RespOK(c.Writer, data, "ok")
}
