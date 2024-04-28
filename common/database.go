package common

import (
	"GinVue/model"
	"fmt"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

// 连接数据库

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginvue"
	username := "root"
	password := "jh529529"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username, password, host, port, database, charset)

	db, err := gorm.Open(driverName, args)

	if err != nil {
		panic("failed to connect databases,err:" + err.Error())
	}

	db.AutoMigrate(&model.User{})
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
