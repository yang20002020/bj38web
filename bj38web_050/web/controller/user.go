package controller

import (
	"bj38web_050/web/utils"
	"fmt"
	"github.com/gin-gonic/gin"
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
	uuid := ctx.Param("uuid")
	fmt.Println("uuid=", uuid)
}
