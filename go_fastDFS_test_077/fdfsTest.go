package main

import (
	"fmt"
	client "github.com/tedcy/fdfs_client"
)

func main() {
	//初始化客户端
	client, err := client.NewClientWithConfig("/etc/fdfs/client.conf")
	if err != nil {
		fmt.Println("初始化客户端错误：", err)
		return
	}

	//上传文件  文件名上传 上传到storage
	resp, err := client.UploadByFilename("zhang.jpg")
	fmt.Println(resp, err)

}
