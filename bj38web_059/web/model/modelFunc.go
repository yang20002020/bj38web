package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 创建全局redis 连接池 句柄
var RedisPool redis.Pool

// 创建函数  初始化连接池
func InitRedis() {
	RedisPool = redis.Pool{
		MaxIdle:         20,
		MaxActive:       50,
		MaxConnLifetime: 60 * 5, //最大生命周期
		IdleTimeout:     60,     //超过60秒就会被断掉数据库
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "192.168.***.***:6379")
		},
	}
}

// 校验图片验证码 该功能需要转移到后端
func CheckImgCode(uuid, imgCode string) bool {
	//链接数据库
	//conn, err := redis.Dial("tcp", "192.168.1.112:6379")
	//if err != nil {
	//	fmt.Println("redis.Dial err:", err)
	//	return false
	//}
	// 从连接池中获取链接
	conn := RedisPool.Get()
	defer conn.Close()

	//查询redis 数据
	//Do(commandName string, args ...interface{}) (reply interface{}, err error)
	code, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println("查询错误  err:", err)
		return false
	}
	//返回校验结果
	return code == imgCode
}

// 存储短信验证码  该功能需要转移到后端
func SaveSmsCode(phone, code string) error {
	//链接redis 连接池 从连接池中获取链接
	conn := RedisPool.Get()
	defer conn.Close()
	//存储图片验证码 到redis中
	_, err := conn.Do("setex", phone+"_code", 60*3, code)
	return err
}
