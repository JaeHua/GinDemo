package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique"`
	Password  string `gorm:"size:255;"`
}

func main() {
	db := initDB()
	defer db.Close()

	r := gin.Default()

	r.POST("/api/auth/register", func(ctx *gin.Context) {

		//获取参数
		name := ctx.PostForm("name")

		telephone := ctx.PostForm("telephone")

		password := ctx.PostForm("password")

		//判断参数合法

		if len(telephone) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "电话必须为11位",
			})
			return
		}
		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "密码不得少于六位",
			})
			return
		}

		if len(name) == 0 {
			name = getRandomName()
		}

		log.Println(name, telephone, password)

		//手机号是否存在
		if isTelephoneExist(db, telephone) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "用户已经存在",
			})
			return
		}

		newUser := User{
			Model:     gorm.Model{},
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)
		//用户验证

		ctx.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
	})
	panic(r.Run())
}

// 返回一个十位的随机字符串
func getRandomName() string {

	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	result := make([]byte, 10)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

// 连接数据库
func initDB() *gorm.DB {
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

	db.AutoMigrate(&User{})
	return db
}

// 判断手机号是否存在
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone=?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
