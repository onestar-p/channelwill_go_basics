/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-05 11:27:37
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-08 10:03:28
 */
package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"channelwill_go_basics/global"
	"channelwill_go_basics/interceptor"
	authpb "channelwill_go_basics/proto/gen/auth/v1"
	etranslatepb "channelwill_go_basics/proto/gen/etranslate/v1"
	"channelwill_go_basics/utils/register/consul"
)

type ClientService struct {
	Name         string
	RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) (err error)
}

func NewClientServices() []*ClientService {
	services := []*ClientService{
		{
			Name:         "auth",
			RegisterFunc: authpb.RegisterAuthServiceHandler, // auth pb
		},
		{
			Name:         "extranslate",                                 // 服务名
			RegisterFunc: etranslatepb.RegisterEtranslateServiceHandler, // etranslate pb
		},
	}
	return services
}

type GRPCConfig struct {
	Name              string             // 服务名
	Addr              string             // 服务地址
	AuthPublicKeyFile string             // 公钥
	RegisterFunc      func(*grpc.Server) // 匿名函数，用户注册服务
}

func RunGRPCServer(c *GRPCConfig) error {

	var (
		appConfig = global.ApplicationConfig
		opts      []grpc.ServerOption // 初始化拦截器
	)

	if c.AuthPublicKeyFile != "" {
		// 实例一个自定义拦截器
		in, err := interceptor.Auth(c.AuthPublicKeyFile)
		if err != nil {
			zap.S().Fatal("cannot create auth interceptor", err.Error())
		}

		// 实例一个grpc，并添加拦截器
		opts = append(opts, grpc.UnaryInterceptor(in))
	}

	s := grpc.NewServer(opts...)

	// 服务注册....
	c.RegisterFunc(s)

	// consul 注册服务
	regsiterClient := consul.NewRegistryClient(appConfig.ConsulInfo.Ip, appConfig.ConsulInfo.Port)
	serviceID := uuid.NewV4() // 随机生成注册服务ID
	err := regsiterClient.Register(
		s,
		appConfig.Ip,
		appConfig.Port,
		appConfig.Name,
		appConfig.ConsulInfo.Tags,
		serviceID.String(),
	)
	if err != nil {
		return fmt.Errorf("service register: %v", err.Error())
	}

	// 监听服务IP地址及端口
	lis, err := net.Listen(appConfig.Network, c.Addr)
	if err != nil {
		return fmt.Errorf("listen : %v", err.Error())
	}

	// 开启协程运行，以便在结束服务时做后续操作
	zap.S().Infof("server started, addr:%s", c.Addr)
	go func() {
		if err := s.Serve(lis); err != nil {
			zap.S().Panic("cannot run gRPC Server", err.Error())
		}
	}()

	// 接收终止信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGALRM)
	<-quit

	// 程序结束时注销 consul 服务
	if err := regsiterClient.DeRegister(serviceID.String()); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")

	// ......

	return nil
}
