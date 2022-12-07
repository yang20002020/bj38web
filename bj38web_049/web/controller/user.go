package controller

import (
	"bj38web_049/web/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 根据gin框架和中文文档业务要求写函数内容
func GetSession(ctx *gin.Context) {
	// 初始化错误返回的map
	resp := make(map[string]string)
	resp["errno"] = utils.RECODE_SESSIONERR
	resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
	//根军gin框架要求 按照json格式进行发送数据 //将数据序列化，返回给浏览器
	ctx.JSON(http.StatusOK, resp)
	fmt.Println("点击按钮登录****")
}
