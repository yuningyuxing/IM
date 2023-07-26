package models

import (
	"context"
	"main/utils"
	"time"
)

// 设置在线用户到redis缓存
// key作为在redis存储用户的键 这里是用户ID
// val表示用户在线信息的字节切片 这里是用户IP地址
// timeTTL表示用户在线信息的过期时间  一但超过这个时间 redis会自动将该信息删除
// time.Duration表示时间间隔的类型 以纳秒为单位表示时间间隔 一般表示持续多少时间后执行某个操作
// 这里四小时后 删除在线用户缓存
func SetUserOnlineInfo(key string, val []byte, timeTTL time.Duration) {
	ctx := context.Background()
	//将用户信息存到redis中
	utils.Red.Set(ctx, key, val, timeTTL)
}
