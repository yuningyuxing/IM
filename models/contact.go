package models

import "gorm.io/gorm"

// 人员关系
type Contact struct {
	gorm.Model
	//谁的关系
	OwnerId uint
	//对应的是谁
	TargetId uint
	//关系类型 0 1 3
	Type int
	//预留字段
	Desc string
}

func (table *Contact) TableName() string {
	return "contact"
}
