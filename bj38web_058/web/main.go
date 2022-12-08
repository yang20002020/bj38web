package main

import (
	"bj38web_058/web/controller"
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

	//router.GET("/api/v1.0/session", controller.GetSession)
	//router.GET("/api/v1.0/imagecode/:uuid", controller.GetImageCd)
	//router.GET("/api/v1.0/smscode/:phone")
	//添加路由分组
	r1 := router.Group("/api/v1.0")
	{
		r1.GET("/session", controller.GetSession)
		//Request URL: http://192.168.63.128:8080/api/v1.0/imagecode/f786f62f-ab9c-4778-9951-a14e165896fc
		r1.GET("/imagecode/:uuid", controller.GetImageCd)
		//http://192.168.1.112:8080/api/v1.0/smscode/15889317897?text=enhe&id=56b89c6b-d62d-45c2-b84d-e9d80d8f187c
		r1.GET("/smscode/:phone", controller.GetSmscd)
	}
	//3启动运行
	router.Run()
}
