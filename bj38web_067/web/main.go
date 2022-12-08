package main

import (
	"bj38web_067/web/controller"
	"bj38web_067/web/model"
	"github.com/gin-gonic/gin"
)

func main() {
	//初始化 全局连接池句柄
	model.InitDb()
	//初始化redis连接池
	model.InitRedis()
	//gin框架开发三步骤
	//1 初始化路由
	router := gin.Default()
	//2 路由匹配
	//router.GET("/", func(context *gin.Context) {
	//	context.Writer.WriteString("项目开始了........")
	//})
	router.Static("/home", "view") //加载页面

	//router.GET("/api/v1.0/session", controller.GetSession)
	//router.GET("/api/v1.0/imagecode/:uuid", controller.GetImageCd)
	//router.GET("/api/v1.0/smscode/:phone")
	//添加路由分组
	r1 := router.Group("/api/v1.0")
	{
		r1.GET("/session", controller.GetSession)
		r1.GET("/imagecode/:uuid", controller.GetImageCd)
		r1.GET("/smscode/:phone", controller.GetSmscd)
		r1.POST("/users", controller.PostRet)
	}
	//3启动运行
	router.Run()
}
