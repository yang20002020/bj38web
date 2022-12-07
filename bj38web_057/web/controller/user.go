package controller

import (
	"bj38web_057/web/proto/getCaptcha"
	"bj38web_057/web/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/gin-gonic/gin"
	"image/color"
	"image/png"
	"net/http"
)

// 根据gin框架和文档业务要求写函数内容
// 获取session信息
func GetSession(ctx *gin.Context) {
	// 初始化错误返回的map
	resp := make(map[string]string)
	resp["errno"] = utils.RECODE_SESSIONERR
	resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
	//根军gin框架要求 按照json格式进行发送数据
	ctx.JSON(http.StatusOK, resp)
}

// 获取图片信息
func GetImageCd(ctx *gin.Context) {
	fmt.Println("*********1111111111*******")
	//获取图片验证码 uuid
	uuid := ctx.Param("uuid")
	fmt.Println("uuid=", uuid)
	/*************************/
	consulReg := consul.NewRegistry()
	consulService := micro.NewService(
		micro.Registry(consulReg),
	)
	//调用 函数 .初始化客户端
	microClient := getCaptcha.NewGetCaptchaService("go.micro.srv.getCaptcha", consulService.Client())
	//调用远程函数
	resp, err := microClient.Call(context.TODO(), &getCaptcha.Request{Uuid: uuid})
	if err != nil {
		fmt.Println("未能找到远程服务:", err)
		return
	}
	//将得到的数据 反序列化，得到图片数据
	var img captcha.Image
	json.Unmarshal(resp.Msg, &img)
	//对数据进行编码  将图片显示到浏览器
	png.Encode(ctx.Writer, img)
}
func ImageCode(ctx *gin.Context) {
	//初始化对象
	cap := captcha.New()
	//设置字体
	cap.SetFont("./conf/comic.ttf") //comic.ttf是一个文件 并且必须在test文件件内
	//设置验证码大小
	cap.SetSize(128, 64)
	//设置干扰强度
	cap.SetDisturbance(captcha.NORMAL)
	//设置前景色
	cap.SetFrontColor(color.RGBA{0, 0, 0, 255})
	//设置背景色
	cap.SetBkgColor(color.RGBA{0, 128, 128, 128}, color.RGBA{255, 255, 10, 255})

	//生成字体
	img, str := cap.Create(6, captcha.ALL)
	fmt.Println("str=", str)
	fmt.Println("*********222222222222*******")
	//对数据进行编码
	png.Encode(ctx.Writer, img)
	//将图片验证码， 结果展示到页面中
	//生成字体  将图片 验证码 展示在页面中
	//http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
	//	img, str := cap.Create(6, captcha.ALL)
	//	png.Encode(w, img)
	//	println(str)
	//})
	////启动服务
	//err := http.ListenAndServe(":8086", nil)
	//if err != nil {
	//	println("err:", err)
	//}
}
