package controller

import (
	"bj38web_062/web/model"
	"bj38web_062/web/proto/getCaptcha"
	userMico "bj38web_062/web/proto/user"
	"bj38web_062/web/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/gin-gonic/gin"
	"image/png"
	"math/rand"
	"net/http"
	"time"
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

// 获取短信验证码
func GetSmscd(ctx *gin.Context) {

	//获取短信验证码
	phone := ctx.Param("phone")
	//http://192.168.****.1***:8080/api/v1.0/smscode/158**********?text=enhe&id=56b89c6b-d62d-45c2-b84d-e9d80d8f187c
	//拆分GET请求中的URL ===格式： 资源路径？/k=v & k=v & k =v
	//参考gin框架 中文文档
	imgCode := ctx.Query("text")
	uuid := ctx.Query("id")
	fmt.Println("---OUT----:", phone, imgCode, uuid)
	/*************需要转移到微服务的代码****************/
	//resp := make(map[string]string)
	////校验图片验证码 是否正确
	//result := model.CheckImgCode(uuid, imgCode)
	//if result {
	//	//校验成功
	//	fmt.Println("校验成功！！！")
	//	//发送短信
	//	err := _main(tea.StringSlice(os.Args[1:]), phone, resp)
	//	if err != nil {
	//		panic(err)
	//	}
	//} else {
	//	//校验失败 显示错误信息
	//	resp["errno"] = utils.RECODE_DATAERR
	//	resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
	//	fmt.Println("校验失败！！！")
	//}
	/******************************************/
	/************微服务添加内容*********/
	//指定consul服务发现
	consulReg := consul.NewRegistry()
	consulService := micro.NewService(
		micro.Registry(consulReg),
	)
	//1.初始化客户端
	microClient := userMico.NewUserService("go.micro.srv.user", consulService.Client())

	//调用远程函数
	resp, err := microClient.SendSms(context.TODO(), &userMico.Request{Phone: phone, ImgCode: imgCode, Uuid: uuid})
	if err != nil {
		fmt.Println("调用远程函数SendSms失败:", err)
		return
	}
	/*********************/
	//发送成功或者失败结果 给浏览器
	ctx.JSON(http.StatusOK, resp)

}

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func _main(args []*string, phone string, resp map[string]string) (_err error) {
	client, _err := CreateClient(tea.String("*************"), tea.String("******************"))
	if _err != nil {
		return _err
	}
	//随机生成一个验证码
	rand.Seed(time.Now().UnixNano())
	//生成六位随机数
	str := rand.Int31n(1000000)
	smsCode := fmt.Sprintf("%06d", str)
	subStr := "{\"code\"" + ":" + "\"" + smsCode + "\"" + "}"
	fmt.Println(subStr)
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("阿里云短信测试"),
		TemplateCode:  tea.String("SMS_1549******"),
		PhoneNumbers:  tea.String(phone), //15********
		TemplateParam: tea.String(subStr),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		// 复制代码运行请自行打印 API 的返回值
		//20221007 发送短信验证码

		_, _err = client.SendSmsWithOptions(sendSmsRequest, runtime)
		if _err != nil {
			//20221007  发送验证码失败
			fmt.Println("发送验证码失败！！！")
			resp["errno"] = utils.RECODE_SMSERR
			resp["errmsg"] = utils.RecodeText(utils.RECODE_SMSERR)
			return _err
		} else {
			//20221007  发送验证码成功
			resp["errno"] = utils.RECODE_OK
			resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
			fmt.Println("发送验证码成功！！！")
			//将短信验证码，存入到redis 数据库
			err := model.SaveSmsCode(phone, smsCode)
			if err != nil {
				fmt.Println("存储短信验证码到redis失败：", err)
				resp["errno"] = utils.RECODE_DBERR
				resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
			}
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	return _err
}

// 发送注册信息
/*{mobile: "15889********", password: "111", sms_code: "94851"}
mobile: "1588*****"
password: "111"
sms_code: "94851" */
func PostRet(ctx *gin.Context) {
	//
	//mobile := ctx.PostForm("mobile")
	//pwd := ctx.PostForm("password")
	//sms_code := ctx.PostForm("sms_code")
	//fmt.Println("mobile=", mobile, "pwd=", pwd, "sms_code=", sms_code)
	//获取数据
	var regData struct {
		Mobile   string "json:\"mobile\""
		PassWord string "json:\"password\""
		SmsCod   string "json:\"sms_code\""
	}
	ctx.Bind(&regData)
	fmt.Println("mobile=", regData.Mobile, "pwd=", regData.PassWord, "sms_code=", regData.SmsCod)
}
