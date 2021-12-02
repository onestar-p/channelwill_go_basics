package initialize

import (
	"context"
	"fmt"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"channelwill_go_basics/global"
	authpb "channelwill_go_basics/proto/auth/gen/v1"
	etranslatepb "channelwill_go_basics/proto/etranslate/gen/v1"
	"channelwill_go_basics/service/auth"
	"channelwill_go_basics/service/etranslate"
	"channelwill_go_basics/utils"
	"channelwill_go_basics/utils/register/servers"
	"channelwill_go_basics/utils/token"
)

var (
	appConfig = global.ApplicationConfig
)

// 网关初始化
func InitGateway() *servers.Servers {
	regServers := servers.NewServers(&servers.ServerConfig{
		ConsulIp:   appConfig.ConsulInfo.Ip,
		ConsulPort: appConfig.ConsulInfo.Port,
		ConsulTags: appConfig.ConsulInfo.Tags,
		AppIp:      appConfig.Ip,
		AppPort:    appConfig.Port,
		AppAddr:    fmt.Sprintf(":%d", appConfig.Port),
	})

	// 网关
	regServers.AddServerRegisterHandler(&servers.ServerRegisterHandlerFunc{
		ServerName: appConfig.AuthSrv,
		RegisterHandlerFunc: func(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
			c := authpb.NewAuthServiceClient(conn)
			if err := authpb.RegisterAuthServiceHandlerClient(ctx, mux, c); err != nil {
				zap.S().Fatal("cannot register auth service handler client", zap.Error(err))
			}
		},
	})

	regServers.AddServerRegisterHandler(&servers.ServerRegisterHandlerFunc{
		ServerName: appConfig.EtranslateSrv,
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
		ConsulIp:   appConfig.ConsulInfo.Ip,
		ConsulPort: appConfig.ConsulInfo.Port,
		ConsulTags: appConfig.ConsulInfo.Tags,
		AppIp:      appConfig.Ip,
		AppPort:    appConfig.Port,
		AppAddr:    fmt.Sprintf(":%d", appConfig.Port),
	})

	// Auth
	regServers.AddServerRegisterServerFunc(&servers.ServerRegisterServerFunc{
		ServerName: appConfig.AuthSrv,
		RegisterServerFunc: func(s *grpc.Server) {

			// 注册auth服务
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				TokenExpire:    appConfig.JwtInfo.Expire * time.Second, // token超时时间
				TokenGenerator: utils.NewToken(token.JWTType),
			})
		},
	})

	// 添加etranslate服务
	regServers.AddServerRegisterServerFunc(&servers.ServerRegisterServerFunc{
		ServerName: appConfig.EtranslateSrv,
		RegisterServerFunc: func(s *grpc.Server) {
			etranslatepb.RegisterEtranslateServiceServer(s, &etranslate.Service{})

		},
	})
	return regServers
}
