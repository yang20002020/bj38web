package model

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 创建全局redis 连接池 句柄
var RedisPool redis.Pool

// 校验图片验证码
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

// 创建函数  初始化连接池
func InitRedis() {
	RedisPool = redis.Pool{
		MaxIdle:         20,
		MaxActive:       50,
		MaxConnLifetime: 60 * 5,
		IdleTimeout:     60, //超过60秒就会被断掉数据库
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "192.168.1.112:6379")
		},
	}
}

// 存储短信验证码
func SaveSmsCode(phone, code string) error {
	//链接redis 连接池 从连接池中获取链接
	conn := RedisPool.Get()
	defer conn.Close()

	//存储图片验证码 到redis中
	_, err := conn.Do("setex", phone+"_code", 60*3, code)
	return err
}

// 校验短信验证码功能
func CheckSmsCode(phone, code string) error {
	//链接redis
	conn := RedisPool.Get()
	//从redis中，根据key 获取value  ---短信验证码码值
	smsCode, err := redis.String(conn.Do("get", phone+"_code"))
	if err != nil {
		fmt.Println("redis get phone_code err:", err)
	}
	//验证码匹配
	if smsCode != code {
		return errors.New("验证码匹配失败！")
	}
	//匹配成功
	return nil
}

// 注册用户信息，写入mysql 数据库
func RegisterUser(mobile, pwd string) error {
	var user User
	user.Name = mobile //默认用手机号作为用户名
	user.Mobile = mobile
	//使用 md5 对pwd 加密
	m5 := md5.New()       //初始md5 对象
	m5.Write([]byte(pwd)) //将密码写入到缓冲区
	//不使用额外的秘钥
	pwd_hash := hex.EncodeToString(m5.Sum(nil))
	user.Password_hash = pwd_hash
	//插入数据到mysql
	return GlobalConn.Create(&user).Error
}
