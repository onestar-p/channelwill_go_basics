/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-04 18:42:16
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-06 10:40:52
 */
package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type RegistryClient interface {
	Register(grpc grpc.ServiceRegistrar, address string, port int, name string, tags []string, id string) error
	DeRegister(serviceID string) error
}

/**
 * @param {string} host Consul 服务器地址
 * @param {int} port Consul 服务器端口
 * @return {*}
 */
func NewRegistryClient(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

type Registry struct {
	Host string
	Port int
}

/**
 * @param {string} address 服务地址
 * @param {int} port 服务端口
 * @param {string} name 服务名
 * @param {[]string} tags 服务标签
 * @param {string} id 服务ID
 * @return {*}
 */
func (r *Registry) Register(
	grpc grpc.ServiceRegistrar,
	address string, port int, name string, tags []string, id string,
) error {

	// 服务注册
	grpc_health_v1.RegisterHealthServer(grpc, health.NewServer())

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	zap.S().Infof("Consul addr: %s", cfg.Address)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Fatalf("annot create Consul: %v", err)
	}

	// 注册中心-生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", address, port),
		Timeout:                        "5s",  // 超时时间
		Interval:                       "5s",  // 检查间隔时间
		DeregisterCriticalServiceAfter: "10s", //
	}

	zap.S().Infof("Consul check addr: %s", check.GRPC)

	// 注册中心-生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Address = address
	registration.Port = port
	registration.Tags = tags
	registration.Check = check
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *Registry) DeRegister(serviceID string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	return client.Agent().ServiceDeregister(serviceID)
}
