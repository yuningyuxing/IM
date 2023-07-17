package models

import (
	"gorm.io/gorm"
)

// 用来描述消息的结构体
type Message struct {
	gorm.Model
	//发送者
	FromId string
	//接收者
	TargetId string
	//消息类型 群聊 私聊 广播
	Type string
	//消息类型 文字 图片 音频
	Media int
	//消息内容
	Content string
	Pic     string
	Url     string
	Desc    string
	//其他数字统计
	Amount int
}

func (table *Message) TableName() string {
	return "message"
}
