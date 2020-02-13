package main

import (
	"fmt"
	"github.com/2020wfw/app"
	"github.com/2020wfw/app/config"
)

func init()  {
	//读取配置文件
	config.LoadConfig()
}

func main(){

	//服务端注册
	//cancel := config.DiscoveryRegister()
	//defer cancel()
	//fmt.Println("服务注册已完成")

	//客户端发现
	config.ResolverWatch()
	fmt.Println("客户端发现已完成")

	//初始化路由表
	router := app.NewRouter()

	//初始化反向代理
	proxy := app.NewProxy(router)
	proxy.StartProxy()

	fmt.Println("网关已启动")

}
