package controller

import (
	"GinVue/common"
	"GinVue/model"
	"GinVue/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func Register(ctx *gin.Context) {

	DB := common.GetDB()
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
		name = util.GetRandomName()
	}

	log.Println(name, telephone, password)

	//手机号是否存在
	if isTelephoneExist(DB, telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "用户已经存在",
		})
		return
	}

	newUser := model.User{
		Model:     gorm.Model{},
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	DB.Create(&newUser)
	//用户验证

	ctx.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
}

// 判断手机号是否存在
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone=?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
