package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 创建中间件
func Test1(ctx *gin.Context) {
	fmt.Println("11111")
}
func main() {
	router := gin.Default()
	//使用中间件
	//type HandlerFunc func(*Context)
	// Use(middleware ...HandlerFunc)
	router.Use(Test1)
	router.GET("/test", func(ctx *gin.Context) {
		fmt.Println("22222")
		ctx.Writer.WriteString("hello  world!")
	})
	router.Run(":9999")
}
