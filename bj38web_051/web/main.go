package main

import (
	"bj38web_051/web/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	//gin框架开发三步骤
	//1 初始化路由
	router := gin.Default()
	//2 路由匹配
	//router.GET("/", func(context *gin.Context) {
	//	context.Writer.WriteString("项目开始了........")
	//})
	router.Static("/home", "view") //加载页面
	router.GET("/api/v1.0/session", controller.GetSession)

	//
	router.GET("/api/v1.0/imagecode/:uuid", controller.GetImageCd)
	//3启动运行
	router.Run()
}
