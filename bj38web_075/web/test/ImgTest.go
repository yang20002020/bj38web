package main

import (
	"github.com/afocus/captcha"
	"image/color"
	"image/png"
	"net/http"
)

func main() {
	//初始化对象
	cap := captcha.New()
	//设置字体
	cap.SetFont("comic.ttf") //comic.ttf是一个文件 并且必须在test文件件内
	//设置验证码大小
	cap.SetSize(128, 64)
	//设置干扰强度
	cap.SetDisturbance(captcha.NORMAL)
	//设置前景色
	cap.SetFrontColor(color.RGBA{0, 0, 0, 255})
	//设置背景色
	cap.SetBkgColor(color.RGBA{0, 128, 128, 128}, color.RGBA{255, 255, 10, 255})
	//生成字体
	//cap.Create(4, captcha.ALL)
	//将图片验证码， 结果展示到页面中
	//生成字体  将图片 验证码 展示在页面中
	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
		img, str := cap.Create(6, captcha.ALL)
		png.Encode(w, img)
		println(str)
	})
	//启动服务
	err := http.ListenAndServe(":8086", nil)
	if err != nil {
		println("err:", err)
	}

}
