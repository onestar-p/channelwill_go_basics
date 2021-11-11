/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-04 16:33:40
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-08 13:55:45
 */
package main

import (
	"context"

	"go.uber.org/zap"

	"channelwill_go_basics/dao/client"
	"channelwill_go_basics/global"
	"channelwill_go_basics/initialize"
	"channelwill_go_basics/utils"
)

func main() {

	// 初始化日志
	initialize.InitLogger()
	// 初始化配置文件
	initialize.InitConfig()

	//动态获取服务端口
	port, err := utils.GetFreePort()
	if err != nil {
		zap.S().Fatalf("cannot get free port: %v", err)
	} else {
		global.ApplicationConfig.Port = port
	}

	// 初始化数据库
	client.InitMysql()
	defer client.Db.Close()

	// 初始化服务
	servers := initialize.InitServers()
	c := context.Background()
	if err := servers.RegisterGRPCServer(c); err != nil {
		zap.S().Panicf("cannot GRPC Run err: %v", err)
	}

}
