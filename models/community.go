package models

import (
	"fmt"
	"gorm.io/gorm"
	"main/utils"
)

// 描述群的
type Community struct {
	gorm.Model
	Name string
	//群拥有者
	OwnerId uint
	//群类型
	Cate int
	//群头像
	Img string
	//群介绍
	Desc string
}

func CreatCommunity(community Community) (int, string) {
	if len(community.Name) == 0 {
		return -1, "群名不能为空"
	}
	if community.OwnerId == 0 {
		return -1, "未登录"
	}
	err := utils.DB.Create(&community).Error
	if err != nil {
		fmt.Println("DB Create community err=", err)
		return -1, "建群失败"
	}
	//建群成功后 我们再加入一下该群
	contact := Contact{}
	contact.OwnerId = community.OwnerId
	contact.TargetId = community.ID
	contact.Type = 2
	utils.DB.Create(&contact)
	return 0, "建群成功"
}

func LoadCommunity(ownerId uint) ([]Community, string) {
	//用来存关系  表示跟这个人有关的群
	contacts := make([]Contact, 0)
	//这个用来存群对应的人ID
	objIds := make([]uint, 0)
	utils.DB.Where("owner_id = ? and type =2", ownerId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, v.TargetId)
	}
	data := make([]Community, 10)
	utils.DB.Where("id in ?", objIds).Find(&data)
	return data, "查询成功"
}

// 加群
func JoinGroup(userId uint, groupName string) (int, string) {
	community := Community{}
	utils.DB.Where("name = ?", groupName).Find(&community)
	if community.Name == "" {
		return -1, "此群不存在"
	} else {
		contact := Contact{}
		contact.OwnerId = userId
		contact.TargetId = community.ID
		contact.Type = 2
		utils.DB.Where("owner_id = ? and target_id = ? and type = 2", userId, community.ID).Find(&contact)
		//因为我们还未创建这个关系 所以当我们查询这个关系在数据库中如果有创建时间 则说明数据库已存在该关系
		if !contact.CreatedAt.IsZero() {
			return -1, "已加入过该群"
		} else {
			utils.DB.Create(&contact)
			return 0, "加群成功"
		}
	}
}
