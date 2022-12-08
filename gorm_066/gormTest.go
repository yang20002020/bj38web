package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

type Student struct {
	ID int
	//string  ---varchar,默认大小255；可以在创建表的时候创建大小
	Name  string `gorm:"size:150;default:\"xiaoming\""`
	Age   int
	Class int       `gorm:"not null"` //非空值
	Join  time.Time `gorm:"type:timestamp"`
}

var GlobaleConn *gorm.DB

func main() {
	//链接数据库  test2 是数据库名字
	//用户名：密码@协议（IP:port）/数据库名
	//链接数据库 获取连接池的句柄 conn
	conn, err := gorm.Open("mysql", "root:fendou2017@tcp(127.0.0.1:3306)/test4?parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("gorm .Open err:", err)
		return
	}
	//defer conn.Close()
	GlobaleConn = conn
	//初始数
	GlobaleConn.DB().SetConnMaxLifetime(10)
	//最大数
	GlobaleConn.DB().SetMaxOpenConns(100)
	//不能使用gorm创建数据库，应该使用sql语句创建好数据库
	// 借助gorm 创建数据表 创建表 students
	//AutoMigrate()默认创建复数形式
	GlobaleConn.SingularTable(true) //使创建的表没有s
	fmt.Println("****************************MAIN****************************")
	fmt.Println(GlobaleConn.AutoMigrate(new(Student)).Error)

	//var stu Student
	//InsertData(stu)
	//删除数据
	//DeleteData()

}
