package main

import (
	"bj38web_082/service/house/handler"
	"bj38web_082/service/house/model"
	house "bj38web_082/service/house/proto/house"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/micro/go-micro/util/log"
)

func main() {
	//初始化mysql 连接池
	model.InitDb()
	//redis 连接池 初始化
	//model.InitRedis()
	//初始化consul对象
	consulReg := consul.NewRegistry()
	// New Service
	service := micro.NewService(
		micro.Address("192.168.63.128:12343"), //防止随机生成port
		micro.Name("go.micro.srv.house"),
		micro.Version("latest"),
		micro.Registry(consulReg),
	)

	// Register Handler
	house.RegisterHouseHandler(service.Server(), new(handler.House))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
