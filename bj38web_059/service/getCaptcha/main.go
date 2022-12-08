package main

import (
	"bj38web_059/service/getCaptcha/handler"
	getCaptcha "bj38web_059/service/getCaptcha/proto/getCaptcha"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/micro/go-micro/util/log"
)

func main() {

	//初始化consul对象
	consulReg := consul.NewRegistry()
	// New Service
	service := micro.NewService(
		micro.Address("192.168.1.112:12341"), //防止随机生成port
		micro.Name("go.micro.srv.getCaptcha"),
		micro.Registry(consulReg),
		micro.Version("latest"),
	)

	// Register Handler
	getCaptcha.RegisterGetCaptchaHandler(service.Server(), new(handler.GetCaptcha))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
