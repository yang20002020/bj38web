package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //代码不直接使用包，底层链接使用
	"github.com/jinzhu/gorm"
)

type Student struct {
	Id   int //成为默认的主键  --主键索引，查询速度快
	Name string
	Age  int
}

// 创建全局连接池句柄
var GloableConn *gorm.DB

func main() {

	//链接数据库  test2 是数据库名字
	//用户名：密码@协议（IP:port）/数据库名
	//链接数据库 获取连接池的句柄 conn ；conn就是一个连接池的句柄
	conn, err := gorm.Open("mysql", "root:fendou2017@tcp(127.0.0.1:3306)/test2")
	if err != nil {
		fmt.Println("gorm .Open err:", err)
		return
	}
	//defer conn.Close()
	GloableConn = conn
	//初始数
	//GloableConn.DB().SetConnMaxLifetime(10)
	//闲置的初始数
	GloableConn.DB().SetMaxIdleConns(10) //idle  闲置的
	//最大数
	GloableConn.DB().SetMaxOpenConns(100)
	//不能使用gorm创建数据库，应该使用sql语句创建好数据库
	// 借助gorm 创建数据表 创建表 students
	//AutoMigrate()默认创建复数形式
	GloableConn.SingularTable(true) //使创建的表没有s
	fmt.Println(GloableConn.AutoMigrate(new(Student)).Error)
}
