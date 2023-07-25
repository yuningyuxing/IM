package models

import (
	"context"
	"main/utils"
	"time"
)

//设置在线用户到redis缓存
func SetUserOnlineInfo(key string, val []byte, timeTTL time.Duration) {
	ctx := context.Background()
	utils.Red.Set(ctx, key, val, timeTTL)
}
