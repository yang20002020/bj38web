package handler

import (
	"context"
	"strconv"

	"github.com/micro/go-micro/util/log"

	"bj38web_081/service/house/model"
	house "bj38web_081/service/house/proto/house"
	"bj38web_081/service/utils"
)

type House struct{}

// 主要功能：将房屋信息写入到数据库中
func (e *House) PubHouse(ctx context.Context, req *house.Request, rsp *house.Response) error {
	log.Log("Received House.Call request")
	//上传房屋业务  把获取到的房屋数据插入数据库
	houseId, err := model.AddHouse(req)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	//{"errno":"0","errmsg":"成功","data":{"house_id":"1"}}
	//house  密码本
	var h house.HouseData //结构体
	// int 类型转换成 string
	h.HouseId = strconv.Itoa(houseId) // h.HouseId 类型 string
	rsp.Data = &h                     //  Data 类型为  *HouseData
	return nil
}

/***************************/
/*
{
	"errno": "0",
	"errmsg": "成功",
	"data": {
	"houses": [
                 {
					"address": "西三旗桥东",
					"area_name": "昌平区",
					"ctime": "2017-11-06 11:16:24",
					"house_id": 1,
					"img_url": "http://101.200.170.171:9998/group1/M00/00/00/Zciqq1oBJY-AL3m8AAS8K2x8TDE052.jpg",
					"order_count": 0,
					"price": 100,
					"room_count": 2,
					"title": "上奥世纪中心",
					"user_avatar": "http://101.200.170.171:9998/group1/M00/00/00/Zciqq1oBLFeALIEjAADexS5wJKs340.png"
                  },
				  {
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
                  }
               ]
}
}*/
////获取已经发布的房源信息
func (e *House) GetHouseInfo(ctx context.Context, req *house.GetReq, resp *house.GetResp) error {
	/*	 返回 参数
	type  GetResp  struct {
		Errno     string
		Errmsg    string
		Data      *GetData
	}
	*/
	//根据用户名获取mysql数据库所有的房屋数据
	//houseInfos 类型 是 []*house.Houses,
	houseInfos, err := model.GetUserHouse(req.UserName)
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	var getData house.GetData //密码本中的数据
	getData.Houses = houseInfos

	resp.Data = &getData
	return nil
}
