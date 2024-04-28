package main

import (
	"GinVue/common"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := common.InitDB() //初始化DB
	defer db.Close()      //数据库延时关闭

	r := gin.Default()
	r = CollectRoute(r) //获取路由

	panic(r.Run())
}
