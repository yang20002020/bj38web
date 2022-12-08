package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 创建中间件
func Test000011(ctx *gin.Context) {
	fmt.Println("11111")

	ctx.Abort()
	fmt.Println("44444")
}

// 创建另外一种中间件
func Test000012() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("33333")
		
		ctx.Next()
		fmt.Println("55555")
	}
}
func main() {
	router := gin.Default()
	//使用中间件
	//type HandlerFunc func(*Context)
	router.Use(Test000011)
	router.Use(Test000012())
	router.GET("/test", func(ctx *gin.Context) {
		fmt.Println("22222")
		ctx.Writer.WriteString("hello  world!")
	})
	router.Run(":9999")
}
