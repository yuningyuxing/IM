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
	//关系类型 1 2 3   1代表好友 2群
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

// 新增好友  这里1添加2时 自动将2添加1
func AddFriend(userId uint, targetName string) (int, string) {
	//通过名字找到要加的那个人
	user := FindUserByName(targetName)
	//不能加自己
	if userId == user.ID {
		return -1, "别加自己"
	}
	contact := Contact{}
	//查询是否关系已经存在 如果存在不能重复添加
	utils.DB.Where("owner_id=? and target_id=? and type = 1", userId, user.ID).Find(&contact)
	if contact.ID != 0 {
		return -1, "无法重复添加"
	}
	//可以正常添加
	if user.Name != "" {
		//开启一个事务 用于确保我插入两条数据时的数据一致性
		tx := utils.DB.Begin()
		contact1 := Contact{}
		contact2 := Contact{}
		contact1.OwnerId = userId
		contact1.TargetId = user.ID
		contact1.Type = 1
		contact2.OwnerId = user.ID
		contact2.TargetId = userId
		contact2.Type = 2

		err := tx.Create(&contact1).Error
		//发送错误 回滚事务
		if err != nil {
			tx.Rollback()
			fmt.Println("Failed to addfriend first", err)
			return -1, "数据库操作出错"
		}
		err = tx.Create(&contact2).Error
		if err != nil {
			tx.Rollback()
			fmt.Println("Failed to addfriend second", err)
			return -1, "数据库操作出错"
		}
		//提交事务 将操作永久化在数据库
		err = tx.Commit().Error
		if err != nil {
			fmt.Println("Failed to addfriend commit tx", err)
			return -1, "数据库出错"
		}
		return 0, "添加成功"
	}
	return -1, "对方不存在"
}

// 通过群 获取群员Id
func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("targetId=? and type =2", communityId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint(v.OwnerId))
	}
	return objIds
}
