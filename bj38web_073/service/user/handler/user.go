package handler

import (
	"bj38web_073/service/user/model"
	user "bj38web_073/service/user/proto/user"
	"bj38web_073/service/utils"
	"context"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"math/rand"
	"os"
	"time"
)

type User struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *User) SendSms(ctx context.Context, req *user.Request, rsp *user.Response) error {

	/*************需要转移到微服务的代码****************/
	resp := make(map[string]string)
	//校验图片验证码 是否正确
	result := model.CheckImgCode(req.Uuid, req.ImgCode)
	if result {
		//校验成功
		fmt.Println("校验成功！！！")
		//发送短信
		err := _main(tea.StringSlice(os.Args[1:]), req.Phone, resp)
		if err != nil {
			panic(err)
		}
	} else {
		//校验失败 显示错误信息
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		fmt.Println("校验失败！！！")
	}
	/******************************************/
	rsp.Errno = resp["errno"]
	rsp.Errmsg = resp["errmsg"]
	return nil
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
	client, _err := CreateClient(tea.String("***************"), tea.String("***************"))
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
		TemplateCode:  tea.String("SMS_154***"),
		PhoneNumbers:  tea.String(phone), //1588***********
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

/************************/
func (e *User) Register(ctx context.Context, req *user.RegReq, rsp *user.Response) error {

	//先校验短信验证码 是否正确 ，在redis中存储短信验证码
	err := model.CheckSmsCode(req.Mobile, req.SmsCode)
	if err == nil {
		//如果校验正确，注册用户、将数据写入到mysql数据库
		err = model.RegisterUser(req.Mobile, req.Password)
		if err != nil {
			rsp.Errno = utils.RECODE_DBERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
			fmt.Println("注册用户信息错误：", err)
			return err
		} else {
			rsp.Errno = utils.RECODE_OK
			rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
		}

	} else {
		//短信验证错误
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		fmt.Println("短信验证码错误:", err)
		return err
	}

	return nil
}
