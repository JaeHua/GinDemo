package controller

import (
	"GinVue/common"
	"GinVue/dto"
	"GinVue/model"
	"GinVue/response"
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
	var requestUser = model.User{}
	ctx.Bind(&requestUser)
	name := requestUser.Name

	telephone := requestUser.Telephone

	password := requestUser.Password

	//判断参数合法

	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "电话必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不得少于六位")

		return
	}

	if len(name) == 0 {
		name = util.GetRandomName()
	}

	//log.Println(name, telephone, password)

	//手机号是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")

		return
	}
	//加密存储
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "加密错误")
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
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		log.Println(err) //记录错误日志
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统错误")

		return

	}
	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func Login(ctx *gin.Context) {
	//获取参数
	DB := common.GetDB()
	//获取参数
	var requestUser = model.User{}
	ctx.Bind(&requestUser)

	telephone := requestUser.Telephone

	password := requestUser.Password

	//验证参数

	log.Println(len(telephone))
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "电话必须为11位")
		return
	}

	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不得少于六位")

		return
	}
	//手机号是否存在
	var user model.User
	DB.Where("telephone=?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	//验证密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Fail(ctx, nil, "密码错误")
		return
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		log.Println(err) //记录错误日志
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统错误")

		return

	}

	response.Success(ctx, gin.H{"token": token}, "登陆成功")
}

// 获取用户信息（已经通过验证，可以从上下文获取信息）

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{"data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
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
