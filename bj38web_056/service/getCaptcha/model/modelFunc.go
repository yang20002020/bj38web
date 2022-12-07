package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 存储图片id 和uuid  键值对 到redis数据库
// 将测试代码redisTest 转移到该函数
func SaveImageCode(code, uuid string) error {
	//1.链接数据库
	conn, err := redis.Dial("tcp", "192.168.1.112:6379")
	if err != nil {
		fmt.Print("redis Dial err:", err)
		return err
	}
	defer conn.Close()
	//fmt.Print("******************11111111111**********")
	//2.操作数据库 有效时间五分钟 key=》 uuid  value=》 code
	_, err = conn.Do("setex", uuid, 60*5, code)
	fmt.Print("******************222222222222222**********")
	return err //如果err 是空返回nil ，如果不是空返回非空值
}
