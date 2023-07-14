package main

//用于测试gorm是否能连通mysql数据库
import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"main/models"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:20020902=QWer@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&models.UserBasic{})

	// Create
	user := &models.UserBasic{}
	user.Name = "yuning"
	db.Create(user)

	// Read
	fmt.Println(db.First(&user, 1)) // 根据整型主键查找

	// Update - 将 product 的 price 更新为 200
	db.Model(&user).Update("Password", "20020902")
	// Update - 更新多个字段
	//db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	// Delete - 删除 product
	//db.Delete(&product, 1)
}
