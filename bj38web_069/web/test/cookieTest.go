package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	//router.SetTrustedProxies([]string{"192.168.1.112"})
	router.GET("/test", func(context *gin.Context) {
		//1.设置cookie maxAge>0 生命周期
		context.SetCookie("mytest", "chuanzhi", 60*60, "", "", false, true)
		//2.maxAge=0 设置临时Cookie ,浏览器关闭 cookie值被干掉
		//context.SetCookie("mytest", "chuanzhi", 0, "", "", false, true)
		context.Writer.WriteString("测试cookie****")
		//3.获取cookie
		cookieVal, _ := context.Cookie("mytest")
		fmt.Println("cookie的value：", cookieVal)

	})
	router.Run(":9999")
}
