/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-04 16:09:55
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-08 09:47:27
 */
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/simplesurance/grpcconsulresolver/consul"
	"google.golang.org/grpc/resolver"

	"channelwill_go_basics/global"
	"channelwill_go_basics/initialize"
)

func init() {
	resolver.Register(consul.NewBuilder())
}

func main() {
	// 初始化日志
	initialize.InitLogger()
	// 初始化配置文件
	initialize.InitConfig()

	appConfig := global.ApplicationConfig

	// 初始化服务
	servers := initialize.InitGateway()
	addr := fmt.Sprintf(":%d", appConfig.HttpPort)
	c := context.Background()
	mux := servers.RegisterGateway(c)
	fmt.Printf("grpc gateway started at %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}

}
