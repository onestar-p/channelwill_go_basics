package servers

import (
	"channelwill_go_basics/interceptor"
	"channelwill_go_basics/utils"
	"channelwill_go_basics/utils/id_generate"
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ServerRegisterHandlerFunc struct {
	ServerName          string
	RegisterHandlerFunc func(c context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux)
}

type ServerRegisterServerFunc struct {
	ServerName         string
	RegisterServerFunc func(*grpc.Server) // 匿名函数，用户注册服务
}

type Servers struct {
	config             *ServerConfig
	RegisterHandler    []*ServerRegisterHandlerFunc // 用于网关注册
	RegisterServerFunc []*ServerRegisterServerFunc  // 匿名函数，用户注册服务
}

type ServerConfig struct {
	ConsulIp   string
	ConsulPort int
	ConsulTags []string
	AppIp      string // 服务IP
	AppAddr    string // 服务地址
	AppPort    int    // 服务端口
}

func NewServers(config *ServerConfig) *Servers {
	return &Servers{
		config: config,
	}
}

func (s *Servers) AddServerRegisterHandler(hander *ServerRegisterHandlerFunc) (err error) {
	s.RegisterHandler = append(s.RegisterHandler, hander)
	return
}

func (s *Servers) AddServerRegisterServerFunc(fun *ServerRegisterServerFunc) (err error) {
	s.RegisterServerFunc = append(s.RegisterServerFunc, fun)
	return
}

// 链接consul
func (s *Servers) ConnConsul(c context.Context, mux *runtime.ServeMux) *runtime.ServeMux {
	for _, handler := range s.RegisterHandler {
		target := fmt.Sprintf("consul://%s:%d/%s?wait=14s", s.config.ConsulIp, s.config.ConsulPort, handler.ServerName)
		zap.S().Infof("Connect Consul: %s -> %s", handler.ServerName, target)
		conn, err := grpc.DialContext(
			c,
			target,
			grpc.WithBlock(),
			grpc.WithInsecure(),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		)
		if err != nil {
			zap.S().Fatalf("cannot connect consul err: %v", err)
		}
		// 注册服务
		handler.RegisterHandlerFunc(c, conn, mux)
	}
	return mux
}

// 注册网关
func (s *Servers) RegisterGateway(c context.Context) *runtime.ServeMux {

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, //解决 protobuf 空值不返回
		&runtime.JSONPb{
			EnumsAsInts: true, // 将枚举的值以整数返回
			OrigName:    true, // 返回的json字段是否以 proto 文件的字段格式返回
		},
	))
	mux = s.ConnConsul(c, mux)
	return mux
}

func (s *Servers) RegisterGRPCServer(c context.Context) error {

	// 监听服务IP地址及端口
	lis, err := net.Listen("tcp", s.config.AppAddr)
	zap.S().Infof("Listen started, addr:%s", s.config.AppAddr)
	if err != nil {
		return fmt.Errorf("listen : %v", err.Error())
	}

	var (
		opts []grpc.ServerOption // 初始化拦截器
	)

	// 注册拦截器
	i := interceptor.NewInterceptor()
	i.AddMiddleware(&interceptor.Auth{})

	// i.AddMiddleware(&interceptor.Mymidd{})

	in, err := i.UnaryServerInterceptor()
	if err != nil {
		zap.S().Fatal("cannot create interceptor", err.Error())
	}
	opts = append(opts, grpc.UnaryInterceptor(in))
	// 注册拦截器 end

	grpcServer := grpc.NewServer(opts...)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	// 服务注册....
	for _, server := range s.RegisterServerFunc {

		server.RegisterServerFunc(grpcServer)
		fmt.Printf("【Server Name: %s】\n", server.ServerName)

		// consul 注册服务
		serviceID, _ := utils.IDGenerate(id_generate.GenUuid).GetID() // 随机生成注册服务ID
		err = s.RegisterConsul(server.ServerName, serviceID)
		if err != nil {
			zap.S().Panic("service register: %v", err.Error())
		}

	}
	zap.S().Infof("server started, addr:%s", s.config.AppAddr)
	if err := grpcServer.Serve(lis); err != nil {
		zap.S().Panic("cannot run gRPC Server", err.Error())
	}
	// 接收终止信号
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGALRM)
	// <-quit

	// // 程序结束时注销 consul 服务
	// if err := regsiterClient.DeRegister(serviceID.String()); err != nil {
	// 	zap.S().Info("注销失败")
	// }
	// zap.S().Info("注销成功")
	return nil
}

func (s *Servers) RegisterConsul(name string, id string) error {
	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", s.config.ConsulIp, s.config.ConsulPort)
	zap.S().Infof("Consul addr: %s", cfg.Address)
	client, err := api.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("annot create Consul: %v", err)
	}

	// 注册中心-生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", s.config.AppIp, s.config.AppPort),
		Timeout:                        "5s",  // 超时时间
		Interval:                       "5s",  // 检查间隔时间
		DeregisterCriticalServiceAfter: "10s", //
	}

	zap.S().Infof("Consul check addr: %s", check.GRPC)
	// 注册中心-生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Address = s.config.AppIp
	registration.Port = s.config.AppPort
	registration.Tags = s.config.ConsulTags
	registration.Check = check
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("cannot ServiceRegister err: %v", err.Error())
	}
	return nil
}
