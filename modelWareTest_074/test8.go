package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// 创建中间件
func Test000111(ctx *gin.Context) {
	fmt.Println("11111")
	t := time.Now()

	ctx.Next()
	fmt.Println(time.Now().Sub(t)) //用当前的事件减去t的时间
}

// 创建另外一种中间件
func Test000112() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("33333")
		//
		//ctx.Abort()
		ctx.Next()
		fmt.Println("55555")
	}
}
func main() {
	router := gin.Default()
	//使用中间件
	//type HandlerFunc func(*Context)
	router.Use(Test000111)
	router.Use(Test000112())
	router.GET("/test", func(ctx *gin.Context) {
		fmt.Println("22222")
		ctx.Writer.WriteString("hello  world!")
	})
	router.Run(":9999")
}
