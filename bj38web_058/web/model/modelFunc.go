package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 校验图片验证码
func CheckImgCode(uuid, imgCode string) bool {
	//链接数据库
	conn, err := redis.Dial("tcp", "192.168.1.112:6379")
	if err != nil {
		fmt.Println("redis.Dial err:", err)
		return false
	}
	defer conn.Close()

	//查询redis 数据
	//Do(commandName string, args ...interface{}) (reply interface{}, err error)
	//code 是string类型
	code, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println("查询错误  err:", err)
		return false
	}
	//返回校验结果
	return code == imgCode
}
