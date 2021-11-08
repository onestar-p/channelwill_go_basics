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

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/simplesurance/grpcconsulresolver/consul"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"

	"channelwill_go_basics/global"
	"channelwill_go_basics/initialize"
	registerServer "channelwill_go_basics/utils/register/server"
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
	addr := fmt.Sprintf(":%d", appConfig.HttpPort)
	c := context.Background()
	conn, _ := grpc.DialContext(
		c,
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", appConfig.ConsulInfo.Ip, appConfig.ConsulInfo.Port, appConfig.Name),
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	// 网关IP地址
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, //解决 protobuf 空值不返回
		&runtime.JSONPb{
			EnumsAsInts: true, // 将枚举的值以整数返回
			OrigName:    true, // 返回的json字段是否以 proto 文件的字段格式返回
		},
	))
	// 注册服务
	for _, s := range registerServer.NewClientServices() {
		if err := s.RegisterFunc(c, mux, conn); err != nil {
			zap.S().Fatalf("cannot register service %s:%s", s.Name, err.Error())
		}
	}

	fmt.Printf("grpc gateway started at %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}

}
