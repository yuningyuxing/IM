package models

//model文件夹用来描述我们要操作的对象
import (
	"fmt"
	"gorm.io/gorm"
	"main/utils"
	"time"
)

// 表示某个用户的信息 用于和数据库关联
type UserBasic struct {
	//自带了一些字段 model
	gorm.Model
	//用户名
	Name string
	//用户密码
	PassWord string
	//用户电话  下面这个是正则表达式
	Phone string `valid;"matches(^1[3-9]{1}\\d{9}$)"`
	//用户邮箱   valid是govalidator库中用来校验邮箱格式是否正确的
	Email string `valid:"email"`
	//用户身份  也就是鉴权 token
	Identity string
	//用户IP
	ClientIp string
	//用户端口
	ClientPort string
	//盐值
	Salt string
	//用户登录时间
	LoginTime uint64
	//用户心跳时间
	HeartbeatTime uint64
	//用户下线时间
	LoginOutTime uint64 `gorm:"colum:login_out_time" json:"login_out_time"`
	//用户是否登录
	IsLogout bool
	//用户设备信息
	DeviceInfo string
}

// 指定在数据库中的表名
func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	return data
}

// 在数据库中创建用户
func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

// 在数据库中删除用户
func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

// 在数据库中更新用户
func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).
		Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email})
}

// 通过名字去寻找用户
func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

// 通过电话去寻找用户
func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

// 通过邮箱去寻找用户
func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
}

// 登录用
func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and pass_word = ?", name, password).First(&user)
	//获取当前系统时间
	str := fmt.Sprintf("%d", time.Now().Unix())
	//进行MD5加密
	temp := utils.MD5Encode(str)
	//然后更新用户的token
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	return user
}
