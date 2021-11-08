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
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	etRoot "channelwill_go_basics"
	"channelwill_go_basics/dao/client"
	"channelwill_go_basics/global"
	"channelwill_go_basics/initialize"
	authpb "channelwill_go_basics/proto/gen/auth/v1"
	etranslatepb "channelwill_go_basics/proto/gen/etranslate/v1"
	"channelwill_go_basics/service/auth"
	"channelwill_go_basics/service/etranslate"
	jwt "channelwill_go_basics/utils/jwt"
	service "channelwill_go_basics/utils/register/server"
)

func main() {

	// 初始化日志
	initialize.InitLogger()
	// 初始化配置文件
	initialize.InitConfig()
	// 初始化数据库
	client.InitMysql()
	defer client.Db.Close()

	var (
		appConfig      = global.ApplicationConfig
		privateKeyFile = etRoot.Path("config/cert/private.key") // 私钥路径
	)

	privKey, err := jwt.NewJWTKey(privateKeyFile).GetPrivateKey()
	if err != nil {
		zap.S().Fatal("cannot pare private key", zap.Error(err))
	}

	var (
		addr          = fmt.Sprintf(":%d", appConfig.Port)
		publicKeyFile = etRoot.Path("config/cert/public.key")
	)
	if err := service.RunGRPCServer(&service.GRPCConfig{
		Name:              appConfig.Name,
		Addr:              addr,
		AuthPublicKeyFile: publicKeyFile,
		RegisterFunc: func(s *grpc.Server) {
			// 注意，在创建 GRPC 服务时，需要确定该服务是否涉及到“加密”：
			// 如有涉及，请传入 “TokenGenerator”。
			// “TokenGenerator”： 用于生成 JWT Token。

			// 注册etranslate服务
			etranslatepb.RegisterEtranslateServiceServer(s, &etranslate.Service{})

			// 注册auth服务
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				TokenExpire:       appConfig.JwtInfo.Expire * time.Second, // token超时时间
				AuthPublicKeyFile: etRoot.Path("config/cert/public.key"),
				TokenGenerator:    jwt.NewJWTTokenGen(appConfig.JwtInfo.Issuer, privKey),
			})
		},
	}); err != nil {
		zap.S().Panicf("cannot GRPC Run err: %v", err)
	}

}
