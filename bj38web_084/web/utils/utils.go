package utils

import (
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
)

// 初始化micro
func InitMicro() micro.Service {
	//初始化客户端
	consulReg := consul.NewRegistry()
	consulService := micro.NewService(
		micro.Registry(consulReg),
	)
	return consulService
}
