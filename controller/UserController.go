package controller

import (
	"GinVue/common"
	"GinVue/model"
	"GinVue/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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

	//log.Println(name, telephone, password)

	//手机号是否存在
	if isTelephoneExist(DB, telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "用户已经存在",
		})
		return
	}
	//加密存储
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "加密错误"})
		return
	}

	newUser := model.User{
		Model:     gorm.Model{},
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)
	//用户验证

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200, "msg": "注册成功"})
}

func Login(ctx *gin.Context) {
	//获取参数
	DB := common.GetDB()

	telephone := ctx.PostForm("telephone")

	password := ctx.PostForm("password")
	//验证参数

	log.Println(len(telephone))
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
	//手机号是否存在
	var user model.User
	DB.Where("telephone=?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "用户不存在",
		})
	}
	//验证密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "密码错误",
		})
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		log.Println(err) //记录错误日志
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统错误"})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200, "data": gin.H{"token": token}, "msg": "登陆成功"})
}

// 获取用户信息（已经通过验证，可以从上下文获取信息）

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"date": gin.H{"user": user},
	})
}

/*
{
    "code": 200,
    "date": {
        "user": {
            "ID": 3,
            "CreatedAt": "2024-04-28T11:35:54Z",
            "UpdatedAt": "2024-04-28T11:35:54Z",
            "DeletedAt": null,
            "Name": "suDVDADhwJ",
            "Telephone": "12345678914",
            "Password": "$2a$10$VzKK3utrIlbBN3JA8NSL9eqgrXt9rKg.3/XswxMq4kyYjcBMQ0A/q"
        }
    }
}
*/

// 判断手机号是否存在
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone=?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
