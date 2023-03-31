package main

import (
	"fmt"
	"go-shop-server/user/global"
	"go-shop-server/user/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	MysqlInit("root", "123456", "127.0.0.1", "3306", "go-shop-server")
}

// mysql连接初始化
func MysqlInit(username, pwd, host, port, dbname string) {
	// dsn := "root:123456@tcp(192.168.0.6:3306)/gin?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		username, pwd, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields: true, //打印sql
		//SkipDefaultTransaction: true, //禁用事务
	})

	if err != nil {
		fmt.Println(err)
	}

	// 自动迁移
	db.AutoMigrate(
		// 表
		model.User{},
	)

	global.DB = db
}
