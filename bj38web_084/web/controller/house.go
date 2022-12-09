package controller

import (
	houseMicro "bj38web_084/web/proto/house"
	"bj38web_084/web/utils"
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

type HouseStu struct {
	Acreage   string   `json:"acreage"`
	Address   string   `json:"address"`
	AreaId    string   `json:"area_id"`
	Beds      string   `json:"beds"`
	Capacity  string   `json:"capacity"`
	Deposit   string   `json:"deposit"`
	Facility  []string `json:"facility"`
	MaxDays   string   `json:"max_days"`
	MinDays   string   `json:"min_days"`
	Price     string   `json:"price"`
	RoomCount string   `json:"room_count"`
	Title     string   `json:"title"`
	Unit      string   `json:"unit"`
}

// 发布房源
/* 浏览器内容
{title: "朗豪坊", price: "1198", area_id: "1", address: "北京大学未名湖", room_count: "4", acreage: "120",…}
acreage:"120"
address:"北京大学未名湖"
area_id:"1"
beds:"双人床"
capacity:"4"
deposit:"800"
facility:
["1", "2", "3", "4", "5", "6", "11", "12", "13", "14", "21", "22", "23"]
max_days:"365"
min_days:"30"
price:"1198"
room_count:"4"
title:"朗豪坊"
unit:"三室一厅"
*/
func PostHouses(ctx *gin.Context) {
	//获取数据   	bind数据的时候不带自动转换
	var house HouseStu
	err := ctx.Bind(&house)

	//校验数据
	if err != nil {
		fmt.Println("获取数据错误", err)
		return
	}

	//获取用户名 //web main.go 程序中已经定义了sessions的key 为userName
	userName := sessions.Default(ctx).Get("userName")

	//处理数据  服务端处理
	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.InitMicro().Client())

	//调用远程服务 传入输入数据
	//resp  返回数据 {"errno":"0","errmsg":"成功","data":{"house_id":"1"}}
	resp, _ := microClient.PubHouse(context.TODO(), &houseMicro.Request{
		Acreage:   house.Acreage,
		Address:   house.Address,
		AreaId:    house.AreaId,
		Beds:      house.Beds,
		Capacity:  house.Capacity,
		Deposit:   house.Deposit,
		Facility:  house.Facility,
		MaxDays:   house.MaxDays,
		MinDays:   house.MinDays,
		Price:     house.Price,
		RoomCount: house.RoomCount,
		Title:     house.Title,
		Unit:      house.Unit,
		UserName:  userName.(string),
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

// 获取已发布房源信息  假数据
func GetUserHouses(ctx *gin.Context) {
	//假数据
	/*
		resp := make(map[string]interface{})

		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
		resp["data"] = ""

		ctx.JSON(http.StatusOK, resp)
	*/
	//获取用户名
	userName := sessions.Default(ctx).Get("userName")

	/*//测试一对多查询
	//有用户名
	var userInfo model.User
	if err := model.GlobalDB.Where("name = ?",userName).Find(&userInfo).Error;err != nil {
		fmt.Println("获取当前用户信息错误",err)
	}
	//房源信息   一对多查询
	var houses []model.House

	model.GlobalDB.Model(&userInfo).Related(&houses)
	fmt.Println("11111111",houses)*/

	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.InitMicro().Client())
	//调用远程服务
	resp, _ := microClient.GetHouseInfo(context.TODO(), &houseMicro.GetReq{UserName: userName.(string)})
	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

// 上传房屋图片
func PostHousesImage(ctx *gin.Context) {
	//获取数据
	houseId := ctx.Param("id")

	fileHeader, err := ctx.FormFile("house_image")

	//校验数据
	if houseId == "" || err != nil {
		fmt.Println("传入数据不完整", err)
		return
	}

	//三种校验 大小,类型,防止重名  fastdfs
	if fileHeader.Size > 50000000 {
		fmt.Println("文件过大,请重新选择")
		return
	}

	fileExt := path.Ext(fileHeader.Filename)
	if fileExt != ".png" && fileExt != ".jpg" {
		fmt.Println("文件类型错误,请重新选择")
		return
	}

	//获取文件字节切片
	file, _ := fileHeader.Open()
	buf := make([]byte, fileHeader.Size)
	file.Read(buf)

	//处理数据  服务中实现
	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.InitMicro().Client())

	//调用远程服务
	resp, _ := microClient.UploadHouseImg(context.TODO(), &houseMicro.ImgReq{
		HouseId: houseId,
		ImgData: buf,
		FileExt: fileExt,
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

// 轮播图
func GetIndex(ctx *gin.Context) {
	// consul micro
	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.InitMicro().Client())

	//调用远程服务  不需要输入参数
	resp, _ := microClient.GetIndexHouse(context.TODO(), &houseMicro.IndexReq{})

	ctx.JSON(http.StatusOK, resp)
}
