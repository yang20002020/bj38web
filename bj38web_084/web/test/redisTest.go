package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func main() {
	//1.链接数据库
	conn, err := redis.Dial("tcp", "192.168.63.128:6379")
	if err != nil {
		fmt.Print("redis Dial err:", err)
		return
	}
	defer conn.Close()
	//fmt.Print("******************11111111111**********")
	//2.操作数据库
	reply, err := conn.Do("set", "itcast11", "itheima")
	fmt.Print("******************222222222222222**********")
	//3.回复助手类函数 ----确定成具体的数据类型
	r, e := redis.String(reply, err)
	fmt.Println("r", r)
	if e != nil {
		fmt.Println("e:", e)
	}
	fmt.Print("******************3333333333333333**********")

}
