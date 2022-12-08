package model

import (
	"crypto/md5"
	"encoding/hex"
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

// 处理登录业务 根据手机号和密码获取用户名
func Login(mobile, pwd string) (string, error) {
	var user User
	//对参数pwd 做md5 hash
	m5 := md5.New()
	m5.Write([]byte(pwd))
	pwd_hash := hex.EncodeToString(m5.Sum(nil))
	//mobile name  都是mysql指定的字段
	// select name from user where  mobile=** and pwd_hash=***
	err := GlobalConn.Where("mobile=?", mobile).Select("name").Where("password_hash=?", pwd_hash).Find(&user).Error
	return user.Name, err
}

// 从mysql中 获取用户信息
func GetUserInfo(userName string) (User, error) {
	// 实现 sql :select *from user where name=username
	var user User
	err := GlobalConn.Where("name=?", userName).First(&user).Error
	return user, err
}
