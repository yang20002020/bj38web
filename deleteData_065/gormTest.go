package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//	type Model struct {
//		ID        uint `gorm:"primary_key"`
//		CreatedAt time.Time
//		UpdatedAt time.Time
//		DeletedAt *time.Time `sql:"index"`
//	}
type Student struct {
	//Id int
	gorm.Model //匿名成员，  这里起继承的作用
	Name       string
	Age        int
}

var GlobaleConn *gorm.DB

func main() {

	//链接数据库  test2 是数据库名字
	//用户名：密码@协议（IP:port）/数据库名
	//链接数据库 获取连接池的句柄 conn
	conn, err := gorm.Open("mysql", "root:fendou2017@tcp(127.0.0.1:3306)/test3?parseTime=True&loc=Local")
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

	var stu Student
	InsertData(stu)
	//删除数据
	DeleteData()
}

// 插入数据
func InsertData(stu Student) {
	fmt.Println("****************************InsertData****************************")
	//创建数据

	stu.Name = "lisi"
	stu.Age = 100

	//插入数据
	fmt.Println(GlobaleConn.Create(&stu).Error)
	fmt.Println(stu)

}

// 通过sql语句 select *from student; 验证
func DeleteData() {
	fmt.Println("****************DeleteData*****************************")
	//软删除  // 通过sql语句 select *from student; 验证
	fmt.Println(GlobaleConn.Where("name=?", "lisi").Delete(new(Student)).Error)
	var stus []Student
	GlobaleConn.Find(&stus)

	fmt.Println(stus) //查询结果为空 但是sql语句查询结果不是空的

	//查询软删除的数据
	fmt.Println("查询软删除数据")
	GlobaleConn.Unscoped().Find(&stus)
	fmt.Println(stus)

	//物理删除  通过 select *from student; 验证
	fmt.Println(GlobaleConn.Unscoped().Where("name=?", "lisi").Delete(new(Student)).Error)

}
