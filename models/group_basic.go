package models

//这个没用到
import "gorm.io/gorm"

// 群信息
type GroupBasic struct {
	gorm.Model
	//群名
	Name string
	//群主
	OwnerId uint
	Icon    string
	Type    int
	Desc    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
