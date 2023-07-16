package utils

//配置一些初始化的地方

//viper库用于处理配置文件的加载和解析
import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

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
	//newLogger 自定义地日志记录器 用于打印SQL语句和数据库操作地日志
	//logger.New是gorm库中地函数 用于创建一个新地日志记录器
	newLogger := logger.New(
		//log.New是Go标准库中的函数 用于创建一个新的日志记录器
		//下面的第二个参数用于设置日志的换行格式 log.LstdFlags是一个log包中的常量 用于设置日志的格式
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢SQL阈值
			LogLevel:      logger.Info, //级别
			Colorful:      true,        //彩色
		},
	)

	//gorm.config是gorm库中的一个配置对象 用于设置数据库连接的配置 其中Logger属性被设置为之前创建的newLogger
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")),
		&gorm.Config{Logger: newLogger})
}

func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
}

// 定义发布消息所使用的键名
const (
	PublishKet = "websocket"
)

// Publish向指定频道发布消息
// ctx是一个上下文对象 channel是频道名 msg是消息内容
func Publish(ctx context.Context, channel string, msg string) error {
	fmt.Println("Publish....  ", msg)
	//使用redis的Publish方法发布消息
	err := Red.Publish(ctx, channel, msg).Err()
	return err
}

// Subscribe订阅指定频道 并返回接受到的消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	//使用redis的PSubscribe方法订阅频道
	sub := Red.PSubscribe(ctx, channel)
	//接受订阅到的消息
	msg, err := sub.ReceiveMessage(ctx)
	fmt.Println("Subscribe.... ", msg.Payload)
	return msg.Payload, err
}
