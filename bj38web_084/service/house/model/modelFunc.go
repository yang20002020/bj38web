package model

import (
	"bj38web_084/service/house/proto/house"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

// 创建全局redis 连接池 句柄
var RedisPool redis.Pool

// 创建函数  初始化连接池
func InitRedis() {
	RedisPool = redis.Pool{
		MaxIdle:         20,
		MaxActive:       50,
		MaxConnLifetime: 60 * 5,
		IdleTimeout:     60, //超过60秒就会被断掉数据库
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "192.168.63.128:6379")
		},
	}
}

// 发布房屋 -- 向 house 表 插入一行数据.
// 整体逻辑 ：向house表中的对应的每一个字段进行赋值操作
// 传入14个数据
func AddHouse(request *house.Request) (int, error) {
	var houseInfo House //数据库 中的变量
	//给house赋值
	houseInfo.Address = request.Address //string 类型
	houseInfo.Title = request.Title
	houseInfo.Unit = request.Unit
	houseInfo.Beds = request.Beds
	//根据userName获取userId
	var user User
	if err := GlobalDB.Where("name = ?", request.UserName).Find(&user).Error; err != nil {
		fmt.Println("查询当前用户失败", err)
		return 0, err
	}

	//sql中一对多插入,只是给外键赋值
	houseInfo.UserId = uint(user.ID)

	//类型转换
	price, _ := strconv.Atoi(request.Price)
	roomCount, _ := strconv.Atoi(request.RoomCount)
	houseInfo.Price = price
	houseInfo.Room_count = roomCount
	houseInfo.Capacity, _ = strconv.Atoi(request.Capacity)
	houseInfo.Deposit, _ = strconv.Atoi(request.Deposit)
	houseInfo.Min_days, _ = strconv.Atoi(request.MinDays)
	houseInfo.Max_days, _ = strconv.Atoi(request.MaxDays)
	houseInfo.Acreage, _ = strconv.Atoi(request.MaxDays)

	//一对多插入
	areaId, _ := strconv.Atoi(request.AreaId)
	houseInfo.AreaId = uint(areaId)

	//request.Facility    所有的家具  房屋 ; Facility类型是  []string
	//type Facility struct {
	//	Id     int      `json:"fid"`     //设施编号
	//	Name   string   `gorm:"size:32"` //设施名字
	//	Houses []*House //都有哪些房屋有此设施  与房屋表进行关联的
	//}
	//切片的遍历操作 Facility类型是  []string
	//[ "1","2","3","4","5","6", "11","12","15"]
	for _, v := range request.Facility {
		id, _ := strconv.Atoi(v)
		fmt.Println("id:", id)
		var fac Facility
		//sql 语句 ？？？
		if err := GlobalDB.Where("id = ?", id).First(&fac).Error; err != nil {
			fmt.Println("家具id错误", err)
			return 0, err
		}
		fmt.Println("fac:", fac) //fac: {1 无线网络 []} fac: {2 热水淋浴 []}

		//查询到了数据
		//Facilities  类型[]*Facility
		houseInfo.Facilities = append(houseInfo.Facilities, &fac)
	}

	if err := GlobalDB.Create(&houseInfo).Error; err != nil {
		fmt.Println("插入房屋信息失败", err)
		return 0, err
	}

	return int(houseInfo.ID), nil
}

// 获取当前用户发布房源
// 返回密码本中的 []*house.Houses
// 整体思路： 将数据库中的数据传给密码本
func GetUserHouse(userName string) ([]*house.Houses, error) {
	var houseInfos []*house.Houses

	//有用户名
	var user User
	if err := GlobalDB.Where("name = ?", userName).Find(&user).Error; err != nil {
		fmt.Println("获取当前用户信息错误", err)
		return nil, err
	}

	//房源信息   一对多查询
	//????
	var houses []House //数据库 House
	GlobalDB.Model(&user).Related(&houses)
	//参考文档
	/*
	 "address": "北清路郑上路",
	        "area_name": "顺义区",
	        "ctime": "2017-11-06 11:38:54",
	        "house_id": 2,
	        "img_url": "http://101.200.170.171:9998/group1/M00/00/00/Zciqq1oBKtmAC8y8AAZcKg5PznU817.jpg",
	        "order_count": 0,
	        "price": 1000,
	        "room_count": 1,
	        "title": "修正大厦302教室",
	        "user_avatar": "http://101.200.170.171:9998/group1/M00/00/00/Zciqq1oBLFeALIEjAADexS5wJKs340.png"
	*/
	//整体思路： 将数据库中的数据传给密码本
	// houses 数据库
	//house.Houses 密码本
	for _, v := range houses {
		var houseInfo house.Houses // house.Houses 密码本
		houseInfo.Title = v.Title
		houseInfo.Address = v.Address
		houseInfo.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		houseInfo.HouseId = int32(v.ID)
		houseInfo.ImgUrl = "http://192.168.63.128:8080/" + v.Index_image_url
		houseInfo.OrderCount = int32(v.Order_count)
		houseInfo.Price = int32(v.Price)
		houseInfo.RoomCount = int32(v.Room_count)
		houseInfo.UserAvatar = "http://192.168.63.128:8080/" + user.Avatar_url

		//获取地域信息
		var area Area
		//related函数可以是以主表关联从表,也可以是以从表关联主表
		GlobalDB.Where("id = ?", v.AreaId).Find(&area)
		houseInfo.AreaName = area.Name
		//密码本数据
		houseInfos = append(houseInfos, &houseInfo)
	}
	return houseInfos, nil
}

// 把图片的凭证存储到数据中 更新 主图,次图  第一张图片是主图,剩下的图片是副图
func SaveHouseImg(houseId, imgPath string) error {
	/*return GlobalDB.Model(new(House)).Where("id = ?",houseId).
	Update("index_image_url",imgPath).Error*/

	//如何判断上传的图是当前房屋的第一张图片
	var houseInfo House //数据库
	//根据 houseId 查找house
	if err := GlobalDB.Where("id = ?", houseId).Find(&houseInfo).Error; err != nil {
		fmt.Println("查询不到房屋信息", err)
		return err
	}

	//说明没有上传过图片  现在上传的图片是主图
	if houseInfo.Index_image_url == "" {
		return GlobalDB.Model(new(House)).Where("id = ?", houseId).
			Update("index_image_url", imgPath).Error
	}

	//上传的幅图
	var houseImg HouseImage
	/*	type HouseImage struct {
				Id      int    `json:"house_image_id"`
				Url     string `gorm:"size:256" json:"url"`
				HouseId uint   `json:"house_id"`
			}
		   其中 Id 是主键  自动添加
	*/

	houseImg.Url = imgPath
	hId, _ := strconv.Atoi(houseId)
	houseImg.HouseId = uint(hId)

	return GlobalDB.Create(&houseImg).Error
}

// 获取房屋信息
func GetIndexHouse() ([]*house.Houses, error) {

	var housesResp []*house.Houses

	var houses []House
	if err := GlobalDB.Limit(5).Find(&houses).Error; err != nil {
		fmt.Println("获取房屋信息失败", err)
		return nil, err
	}

	for _, v := range houses {
		var houseTemp house.Houses
		houseTemp.Address = v.Address
		//根据房屋信息获取地域信息
		var area Area
		var user User

		GlobalDB.Model(&v).Related(&area).Related(&user)

		houseTemp.AreaName = area.Name
		houseTemp.Ctime = v.CreatedAt.Format("2006-01-02 15:04:05")
		houseTemp.HouseId = int32(v.ID)
		houseTemp.ImgUrl = "http://192.168.63.128:8888/" + v.Index_image_url
		houseTemp.OrderCount = int32(v.Order_count)
		houseTemp.Price = int32(v.Price)
		houseTemp.RoomCount = int32(v.Room_count)
		houseTemp.Title = v.Title
		houseTemp.UserAvatar = "http://192.168.63.128:8888/" + user.Avatar_url

		housesResp = append(housesResp, &houseTemp)
	}

	return housesResp, nil
}
