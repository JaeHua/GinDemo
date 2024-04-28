package main

import (
	"GinVue/common"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"os"
)

func main() {
	InitConfig()

	db := common.InitDB() //初始化DB
	defer db.Close()      //数据库延时关闭

	r := gin.Default()
	r = CollectRoute(r) //获取路由

	port := viper.GetString("server.port")

	panic(r.Run(":" + port))
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
