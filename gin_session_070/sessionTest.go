package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	//"github.com/gin-contrib/sessions"
)

func main() {
	router := gin.Default()
	//初始化容器
	store, _ := redis.NewStore(10, "tcp", "192.168.*****.****:6379", "", []byte("secret"))
	//使用容器 cookie 名称为mysession
	//设置临时session
	//store.Options(sessions.Options{
	//	MaxAge: 0,
	//})
	//mysession 是cookie的值
	router.Use(sessions.Sessions("mysession", store))
	router.GET("/test", func(context *gin.Context) {
		//调用session
		s := sessions.Default(context)
		//设置session
		s.Set("itcast", "itheima")
		//修改session时，需要save函数配合，否则不生效
		s.Save()
		context.Writer.WriteString("测试 session****")

		//获取session
		v := s.Get("itcast")
		fmt.Println("获取session：", v.(string))
	})
	router.Run(":9999")
}
