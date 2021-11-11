package initialize

import (
	etRoot "channelwill_go_basics"
	"channelwill_go_basics/global"
	authpb "channelwill_go_basics/proto/auth/gen/v1"
	etranslatepb "channelwill_go_basics/proto/etranslate/gen/v1"
	"channelwill_go_basics/service/auth"
	"channelwill_go_basics/service/etranslate"
	"channelwill_go_basics/utils/jwt"
	"channelwill_go_basics/utils/register/servers"
	"context"
	"fmt"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	appConfig     = global.ApplicationConfig
	publicKeyFile = etRoot.Path("config/cert/public.key") // 公钥路径
)

// 网关初始化
func InitGateway() *servers.Servers {
	regServers := servers.NewServers(&servers.ServerConfig{
		ConsulIp:      appConfig.ConsulInfo.Ip,
		ConsulPort:    appConfig.ConsulInfo.Port,
		ConsulTags:    appConfig.ConsulInfo.Tags,
		AppIp:         appConfig.Ip,
		AppPort:       appConfig.Port,
		AppAddr:       fmt.Sprintf(":%d", appConfig.Port),
		PublicKeyFile: publicKeyFile,
	})

	// 网关
	regServers.AddServerRegisterHandler(&servers.ServerRegisterHandlerFunc{
		ServerName: "auth",
		RegisterHandlerFunc: func(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
			c := authpb.NewAuthServiceClient(conn)
			if err := authpb.RegisterAuthServiceHandlerClient(ctx, mux, c); err != nil {
				zap.S().Fatal("cannot register auth service handler client", zap.Error(err))
			}
		},
	})

	regServers.AddServerRegisterHandler(&servers.ServerRegisterHandlerFunc{
		ServerName: "etranslate",
		RegisterHandlerFunc: func(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
			c := etranslatepb.NewEtranslateServiceClient(conn)
			if err := etranslatepb.RegisterEtranslateServiceHandlerClient(ctx, mux, c); err != nil {
				zap.S().Fatal("cannot register auth service handler client", zap.Error(err))
			}
		},
	})
	return regServers
}

// 初始化服务
func InitServers() *servers.Servers {
	regServers := servers.NewServers(&servers.ServerConfig{
		ConsulIp:      appConfig.ConsulInfo.Ip,
		ConsulPort:    appConfig.ConsulInfo.Port,
		ConsulTags:    appConfig.ConsulInfo.Tags,
		AppIp:         appConfig.Ip,
		AppPort:       appConfig.Port,
		AppAddr:       fmt.Sprintf(":%d", appConfig.Port),
		PublicKeyFile: publicKeyFile,
	})

	// Auth
	regServers.AddServerRegisterServerFunc(&servers.ServerRegisterServerFunc{
		ServerName: "auth",
		RegisterServerFunc: func(s *grpc.Server) {

			privateKeyFile := etRoot.Path("config/cert/private.key") // 私钥路径
			privKey, err := jwt.NewJWTKey(privateKeyFile).GetPrivateKey()
			if err != nil {
				zap.S().Fatal("cannot pare private key", zap.Error(err))
			}
			// 注册auth服务
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				TokenExpire:       appConfig.JwtInfo.Expire * time.Second, // token超时时间
				AuthPublicKeyFile: publicKeyFile,
				TokenGenerator:    jwt.NewJWTTokenGen(appConfig.JwtInfo.Issuer, privKey),
			})
		},
	})

	// 添加etranslate服务
	regServers.AddServerRegisterServerFunc(&servers.ServerRegisterServerFunc{
		ServerName: "etranslate",
		RegisterServerFunc: func(s *grpc.Server) {
			etranslatepb.RegisterEtranslateServiceServer(s, &etranslate.Service{})

		},
	})
	return regServers
}
