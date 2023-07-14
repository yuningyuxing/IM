package utils

//viper库用于处理配置文件的加载和解析
import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

//配置一些初始化的地方

// 初始化去取得config文件夹中的app.yml文件中存储的数据库配置
func InitConfig() {
	//指定要读取的配置文件的名称
	viper.SetConfigName("app")
	//用于添加配置文件的搜索路径 可以多次调用 以添加多个搜索路径
	viper.AddConfigPath("config")
	//用于读取和解析配置文件 解析后的配置数据将存储在viper的内部结构中
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config mysql", viper.Get("mysql"))
}

// 初始化数据库
func InitMysql() {
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{})
}
