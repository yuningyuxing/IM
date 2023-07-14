package models

//model文件夹用来描述我们要操作的对象
import "gorm.io/gorm"

// 表示某个用户的信息 用于和数据库关联
type UserBasic struct {
	//自带了一些字段 model
	gorm.Model
	//用户名
	Name string
	//用户密码
	PassWord string
	//用户电话
	Phone string
	//用户邮箱
	Email string
	//用户身份?
	Identity string
	//用户IP
	ClientIp string
	//用户端口
	ClientPort string
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

// 给实体绑定一个方法
func (table *UserBasic) TableName() string {
	return "user_basic"
}
