package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 创建中间件
func Test01(ctx *gin.Context) {
	fmt.Println("11111")
}

// 创建另外一种中间件
func Test02() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("33333")
	}
}
func main() {
	router := gin.Default()
	//使用中间件
	//type HandlerFunc func(*Context)
	router.Use(Test01)
	router.Use(Test02())
	router.GET("/test", func(ctx *gin.Context) {
		fmt.Println("22222")
		ctx.Writer.WriteString("hello  world!")
	})
	router.Run(":9999")
}
