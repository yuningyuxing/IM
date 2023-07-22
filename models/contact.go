package models

import (
	"fmt"
	"gorm.io/gorm"
	"main/utils"
)

// 人员关系
type Contact struct {
	gorm.Model
	//谁的关系
	OwnerId uint
	//对应的是谁
	TargetId uint
	//关系类型 0 1 3   1代表好友
	Type int
	//预留字段
	Desc string
}

func (table *Contact) TableName() string {
	return "contact"
}

// 获取好友列表
func SearchFriends(userId uint) []UserBasic {
	//这个用来存所有的我当前这个人的好友关系的结构体contact
	contacts := make([]Contact, 0)
	//用来存跟我为好友的ID
	objIds := make([]uint64, 0)
	//搜索我的关系中 是好友关系的contact
	utils.DB.Where("owner_id = ? and type = 1", userId).Find(&contacts)
	//把这些关系中的好友ID提取出来
	for _, v := range contacts {
		fmt.Println(v)
		//这里存入好友的ID
		objIds = append(objIds, uint64(v.TargetId))
	}
	//用来存所有的好友
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}
