package handler

import (
	getCaptcha "bj38web_053/service/getCaptcha/proto/getCaptcha"
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"image/color"
)

type GetCaptcha struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetCaptcha) Call(ctx context.Context, req *getCaptcha.Request, rsp *getCaptcha.Response) error {
	// 字节数组
	imgBuf := ImageCode()
	//将 imgBuf 使用 参数 rsp 传出

	rsp.Msg = imgBuf
	return nil
}

// 从conroller/user.go 中剪切代码的时候，注意gin 属于前端框架，后端不需要，//获取图片验证码 uuid相关代码不需要
func ImageCode() []byte {
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

	//对数据进行编码 此处代码屏蔽
	//png.Encode(ctx.Writer, img)

	//将生成的图片 序列化 得到json文件
	imgBuf, err := json.Marshal(img)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return imgBuf

}
