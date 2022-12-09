package main

import (
	"bj38web_084/web/controller"
	"bj38web_084/web/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func LoginFileter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := sessions.Default(ctx)
		userName := s.Get("userName")
		if userName == nil {
			//如果结果为空，就没有必要进入下一个中间件
			ctx.Abort() //从这里返回，不必继续执行
		} else {
			ctx.Next() //继续向下
		}
	}
}
func main() {
	//初始化 全局连接池句柄
	model.InitDb()
	//初始化redis连接池
	model.InitRedis()
	//初始化路由
	router := gin.Default()
	//初始化容器
	store, _ := redis.NewStore(10, "tcp", "192.168.63.128:6379", "", []byte("secret"))
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
		r1.GET("/session", controller.GetSession)         //user
		r1.GET("/imagecode/:uuid", controller.GetImageCd) //user
		r1.GET("/smscode/:phone", controller.GetSmscd)    //user
		r1.POST("/users", controller.PostRet)             //user
		r1.GET("/areas", controller.GetArea)              //user
		r1.POST("/sessions", controller.PostLogin)        //post之后才有session
		// 添加中间件，以后的路由不用再校验session了
		r1.Use(LoginFileter())
		r1.DELETE("/session", controller.DeleteSession) //user
		r1.GET("/user", controller.GetUserInfo)         //user
		r1.PUT("/user/name", controller.PutUerInfo)     //user 更新用户名
		r1.POST("/user/avatar", controller.PostAvatar)  //user  上传头像
		// 上传实名认证
		r1.POST("/user/auth", controller.PostUserInfo) // user 上传实名认证
		//获取实名认证 ？？？
		r1.GET("/user/auth", controller.GetUserInfo) // user  获取用户信息
		//发布房源 点击“发布房源信息”
		r1.POST("/houses", controller.PostHouses) // house 发布房源信息
		//获取已经发布的房源信息 点击“我的房源”
		r1.GET("/user/houses", controller.GetUserHouses) // house  获取已经发布的房源信息
		//添加房源图片 点击“上传”
		r1.POST("/houses/:id/images", controller.PostHousesImage) //house 添加图片  上传
		//展示首页轮播图  ???
		r1.GET("/house/index", controller.GetIndex) // house 展示 首页 轮播图

	}
	//3启动运行
	router.Run()
}
