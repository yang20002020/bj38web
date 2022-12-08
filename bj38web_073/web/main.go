package main

import (
	"bj38web_073/web/controller"
	"bj38web_073/web/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	//初始化 全局连接池句柄
	model.InitDb()
	//初始化redis连接池
	model.InitRedis()
	//初始化路由
	router := gin.Default()
	//初始化容器
	store, _ := redis.NewStore(10, "tcp", "192.168.1.112:6379", "", []byte("secret"))
	//使用容器 store
	router.Use(sessions.Sessions("mysession", store))
	//gin框架开发三步骤
	//1 初始化路由
	//router := gin.Default()
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
		r1.GET("/areas", controller.GetArea)
		r1.POST("/sessions", controller.PostLogin)
		r1.DELETE("/session", controller.DeleteSession)
		r1.GET("/user", controller.GetUserInfo)
		r1.PUT("/user/name", controller.PutUerInfo)
	}
	//3启动运行
	router.Run()
}
