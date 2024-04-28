package middleware

import (
	"GinVue/common"
	"GinVue/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//获取Authorization header
		tokenString := ctx.GetHeader("Authorization")

		//validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort() // 抛弃请求
			return
		}
		//除去Bearer
		tokenString = tokenString[7:]

		token, claim, err := common.ParseToken(tokenString)

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
		}
		//验证通过后获取其中的userID
		userId := claim.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		//用户被消除了
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
		}
		//用户存在,将user写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}
