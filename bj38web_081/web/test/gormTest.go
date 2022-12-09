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

var GlobaleConn *gorm.DB
var stu Student
var stus []Student

func main() {

	//链接数据库  test2 是数据库名字
	//用户名：密码@协议（IP:port）/数据库名
	//链接数据库 获取连接池的句柄 conn
	conn, err := gorm.Open("mysql", "root:fendou2017@tcp(127.0.0.1:3306)/test2")
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

	InsertData(stu)
	SearchData(stu)
	InsertData_1(stu)
	//
	SearchData_1(stu)
	//
	UpdateData(stu)
	delete()
}

// 插入数据
func InsertData(stu Student) {
	fmt.Println("****************************InsertData****************************")
	//创建数据

	stu.Name = "zhangsan"
	stu.Age = 100

	//插入数据
	fmt.Println(GlobaleConn.Create(&stu).Error)
	fmt.Println(stu)

}
func InsertData_1(stu Student) {
	fmt.Println("****************************InsertData_1****************************")
	//创建数据
	stu.Name = "lisi"
	stu.Age = 90

	//插入数据
	fmt.Println(GlobaleConn.Create(&stu).Error)
	fmt.Println(stu)
}

// 查询数据
func SearchData(stu Student) {
	fmt.Println("****************************SearchData****************************")
	//获取表中的第一个数据 ,按主键排序 查询到第一条 name和age
	//相当于sql语句 select *from student order by id limit 1;
	GlobaleConn.First(&stu)
	fmt.Println(stu)
	fmt.Println("stus:", stus) //切片初始值为0
}
func SearchData_1(stu Student) {
	fmt.Println("****************************SearchData_1****************************")
	fmt.Println(stu)
	//只查询name 和age   不查询其他值
	GlobaleConn.Select("name,age").First(&stu)
	fmt.Println(stu)
	fmt.Println("******************************")
	InsertData_1(stu)
	//获取表中所有数据
	GlobaleConn.Select("name,age").Find(&stu)
	fmt.Println(stu) //??
	fmt.Println("1111111111111111111111111")
	//修改为切片 才能获取所有数据
	//var stus []Student
	//相当于 select *from name,age from student;
	GlobaleConn.Select("name,age").Find(&stus)
	fmt.Println("stus:", stus)
	fmt.Println("22222222222222222222222222222222")
	//获取最后一条记录 按主键排序
	//相当于SQL 语句select  *from  Student order by id desc limit 1;
	GlobaleConn.Select("name,age").Last(&stu)
	fmt.Println(stu)
	fmt.Println("23333333333333333333333333")
	//相当于select name,age from student where name="lisi";
	GlobaleConn.Select("name,age").Where("name=?", "lisi").Find(&stus)
	fmt.Println("stus:", stus)
	fmt.Println("4444444444444444444444444444444444")
	//相当于select name,age from student where name="lisi" and age=20;
	//方法一
	GlobaleConn.Select("name,age").Where("name=?", "lisi").
		Where("age=?", 90).Find(&stus)
	fmt.Println("stus:", stus)
	//方法二
	GlobaleConn.Select("name,age").Where("name=? and age =?", "lisi", 90).Find(&stus)
	fmt.Println("stus:", stus)
	// select where  和 Find 返回值都是DB*  属于链式调用
}
func UpdateData(stu Student) {
	fmt.Println("***********************UpdateData********************")
	//Module 指定需要操作的“student” 表
	//Where("name=?", "lisi")  指定过滤条件
	//Update("name", "zhaoliu") 指定name lisi 更新为zhaoliu
	//update
	fmt.Println("***********************Update********************")
	fmt.Println(GlobaleConn.Model(new(Student)).Where("name=?", "lisi").Update("name", "zhaoliu").Error)
	fmt.Println("stus:", stus)
	GlobaleConn.Select("name,age").Find(&stus)
	fmt.Println("stus:", stus)
	fmt.Println("***********************Updates********************")
	//updates
	fmt.Println(GlobaleConn.Model(new(Student)).Where("name=?", "zhaoliu").
		Updates(map[string]interface{}{"name": "liuqi", "age": 77}).Error)
	fmt.Println("stus:", stus)
	GlobaleConn.Select("name,age").Find(&stus)
	fmt.Println("stus:", stus)

}

func delete() {
	GlobaleConn.Delete(&stu)
	fmt.Println(stu)
}
