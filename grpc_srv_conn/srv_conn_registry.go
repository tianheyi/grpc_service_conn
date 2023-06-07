package grpc_srv_conn

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important, 初始化了解析consul://的resolver
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"grpc_service_conn/service_registry"
)

type SrvManager interface {
	GetGrpcConn(serviceName string, waitSecond int) (*grpc.ClientConn, error)
}

type SrvManagerImpl struct {
	serviceRegistry   service_registry.ServiceRegistry
	grpcServiceConfig string
}

// GrpcServiceConfig grpc服务配置，具体参考：
//   - https://github.com/grpc/grpc/blob/master/doc/service_config.md
//   - https://github.com/grpc/grpc-proto/blob/master/grpc/service_config/service_config.proto
type GrpcServiceConfig struct {
	// https://github.com/grpc/grpc/blob/master/doc/load-balancing.md
	// loadBalancer模块负载均衡策略配置，此字段废弃，只有未配置load_balanding_config的时候才会用此
	LoadBalancingPolicy string `json:"loadBalancingPolicy,omitempty"`
	MethodConfig        []struct {
		Name []struct {
			Service string `json:"service,omitempty"`
		} `json:"name"`
		WaitForReady bool `json:"waitForReady,omitempty"`
		// 重试配置，可与grpc-middleware中的retry拦截器一起使用
		// 即客户端重试不成功之后，返回给retry拦截器继续重试
		RetryPolicy struct {
			MaxAttempts          int      `json:"MaxAttempts,omitempty"`
			InitialBackoff       string   `json:"InitialBackoff,omitempty"`
			MaxBackoff           string   `json:"MaxBackoff,omitempty"`
			BackoffMultiplier    float64  `json:"BackoffMultiplier,omitempty"`
			RetryableStatusCodes []string `json:"RetryableStatusCodes,omitempty"`
		} `json:"retryPolicy,omitempty"`
	} `json:"methodConfig,omitempty"`
}

// 从consul根据服务名称然后根据负载均衡算法获取一个服务信息，返回其grpc连接
// 用完记得调用Close方法断开连接
func (s *SrvManagerImpl) GetGrpcConn(serviceName string, waitSecond int) (*grpc.ClientConn, error) {

	config := s.serviceRegistry.GetConfig()
	var opts []grpc.DialOption
	opts = append(opts,
		// 禁用 SSL/TLS 安全传输协议，使用不安全的明文传输方式。
		grpc.WithInsecure(),
		// 默认配置信息，可以设置负载均衡策略等选项，这里设置为轮询（round_robin）负载均衡。
		//grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //轮询负载均衡 json字符串
		grpc.WithDefaultServiceConfig(s.grpcServiceConfig),
	)
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf(
			// 在查询服务列表时，如果没有找到任何可用的服务实例，就会等待 waitSecond 秒，直到有新的服务实例注册成功
			// 如果为0或者小于0，如果没找到，会立即返回错误
			//
			"consul://%s:%d/%s?wait=%ds",
			config.Host,
			config.Port,
			serviceName,
			waitSecond,
		),
		opts...,
	)
	return conn, err
}

func NewSrvManagerImpl(registry service_registry.ServiceRegistry, config *GrpcServiceConfig) (SrvManager, error) {
	if registry == nil {
		return nil, errors.New("registry is nil.")
	}
	resByte, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	return &SrvManagerImpl{
		serviceRegistry:   registry,
		grpcServiceConfig: string(resByte),
	}, nil
}
