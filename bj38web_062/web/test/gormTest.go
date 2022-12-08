package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //代码不直接使用包，底层链接使用  链接数据库，首先要导入驱动
	"github.com/jinzhu/gorm"
)

type Student struct {
	Id   int //成为默认的主键  --主键索引，查询速度快
	Name string
	Age  int
}

func main() {

	//链接数据库  test2 是数据库名字
	//用户名：密码@协议（IP:port）/数据库名
	conn, err := gorm.Open("mysql", "root:fendou2017@tcp(127.0.0.1:3306)/test2")
	if err != nil {
		fmt.Println("gorm .Open err:", err)
		return
	}
	defer conn.Close()
	//不能使用gorm创建数据库，应该使用sql语句创建好数据库
	// 借助gorm 创建数据表 创建表 students
	//AutoMigrate()默认创建复数形式 自动添加s
	conn.SingularTable(true) //使创建的表没有s
	fmt.Println(conn.AutoMigrate(new(Student)).Error)
}
