package handler

import (
	"bj38web_073/service/getCaptcha/model"
	getCaptcha "bj38web_073/service/getCaptcha/proto/getCaptcha"
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
	imgBuf, str := ImageCode()
	//将 imgBuf 使用 参数 rsp 传出

	//存储图片验证码到redis中
	err := model.SaveImageCode(str, req.Uuid)
	if err != nil {
		return err
	}
	rsp.Msg = imgBuf
	return nil
}

func ImageCode() ([]byte, string) {
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
	img, str := cap.Create(4, captcha.LOWER)

	fmt.Println("str=", str)
	fmt.Println("*********222222222222*******")

	//对数据进行编码 此处代码屏蔽
	//png.Encode(ctx.Writer, img)

	//将生成的图片 序列化 得到json文件
	imgBuf, err := json.Marshal(img)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return imgBuf, str

}
