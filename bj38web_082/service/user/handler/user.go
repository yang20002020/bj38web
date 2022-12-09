package handler

import (
	"bj38web_082/service/user/model"
	microUser "bj38web_082/service/user/proto/user"
	"bj38web_082/service/utils"
	"context"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type User struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *User) SendSms(ctx context.Context, req *microUser.Request, rsp *microUser.Response) error {

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
	client, _err := CreateClient(tea.String("*************"), tea.String("****************"))
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
		TemplateCode:  tea.String("SMS_1******"),
		PhoneNumbers:  tea.String(phone), //15***********
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
func (e *User) Register(ctx context.Context, req *microUser.RegReq, rsp *microUser.Response) error {

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

/***************实名认证*********************/
func (e *User) AuthUpdate(ctx context.Context, req *microUser.AuthReq, resp *microUser.AuthResp) error {

	//调用借口校验realName和idcard是否匹配

	//存储真实姓名和真实身份证号  数据库
	err := model.SaveRealName(req.UserName, req.RealName, req.IdCard)
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	return nil
}

/************获取用户信息*******************/
func (e *User) GetUsrInfo(ctx context.Context, req *microUser.UserName, resp *microUser.UserInfoResp) error {

	userName := req.UserName
	/**********以下部分内容已经转移到服务端***************/
	//根据用户名获取用户信息 ---查mysql数据库
	user, err := model.GetUserInfo(userName)
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	//{"data":{"avatar_url":"http://192.168.*****.*****:8080/group1/M00/00/00/wKg_gGOAZq6Abq7mAADe1y314UQ299.png",
	//         "id_card":"13120151224",
	//         "mobile":"15*****",
	//         "name":"15*****",
	//         "real_name":"yangzhang",
	//         "user_id":28},
	//  "errmsg":"成功","errno":"0"}
	//参考 浏览器 user 中的回应response，获得对应的字段 或者参考文档
	//temp := make(map[string]interface{})
	//temp["user_id"] = user.ID
	//temp["name"] = user.Name
	//temp["mobile"] = user.Mobile
	//temp["real_name"] = user.Real_name
	//temp["id_card"] = user.Id_card
	//temp["avatar_url"] = "http://192.168.63.128:8080/" + user.Avatar_url

	//resp["data"] = temp

	//定义密码本中变量
	var userInfo microUser.UserInfo
	userInfo.UserId = strconv.Itoa(user.ID)
	userInfo.Name = user.Name
	userInfo.Mobile = user.Mobile
	userInfo.RealName = user.Real_name
	userInfo.IdCard = user.Id_card
	userInfo.AvatarUrl = "http://192.168.****.****:8080/" + user.Avatar_url
	resp.Data = &userInfo
	/********************************************/

	return nil
}
