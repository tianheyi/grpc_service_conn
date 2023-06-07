package service_registry

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"time"
)

type ConsulServiceRegistryImpl struct {
	client *api.Client
	config *ServiceRegistryConfig
}

func (c *ConsulServiceRegistryImpl) Register(id string, address string, port int, name string, tags []string) error {
	// 生成健康检查对象
	checkObj := &api.AgentServiceCheck{
		//srv层用GRPC web层用HTTP
		HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	// 生成注册对象
	registrationObj := &api.AgentServiceRegistration{
		ID:      id,
		Address: address,
		Port:    port,
		Name:    name,
		Tags:    tags,
		Check:   checkObj,
	}
	// 注册
	err := c.client.Agent().ServiceRegister(registrationObj)
	return err
}

func (c *ConsulServiceRegistryImpl) DeRegister(serviceId string) error {
	err := c.client.Agent().ServiceDeregister(serviceId)
	return err
}

func (c *ConsulServiceRegistryImpl) GetServiceList(service string, tags []string, passingOnly bool, handler func([]*ServiceNode) *ServiceNode) ([]*ServiceNode, error) {
	serviceList, _, err := c.client.Health().ServiceMultipleTags(service, tags, passingOnly, nil)
	if err != nil {
		return nil, err
	}
	// 转换
	nodeList := make([]*ServiceNode, len(serviceList))
	for i, serviceNode := range serviceList {
		nodeList[i] = &ServiceNode{
			Address: serviceNode.Service.Address,
			Port:    serviceNode.Service.Port,
		}
	}
	if handler == nil {
		return nodeList, err
	} else {
		nodeList = []*ServiceNode{handler(nodeList)}
		return nodeList, nil
	}
}

func (c *ConsulServiceRegistryImpl) GetConfig() *ServiceRegistryConfig {
	return c.config
}

func NewCustomConsulClient(config *ServiceRegistryConfig) ServiceRegistry {
	// 配置、连接consul
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", config.Host, config.Port)
	//// todo
	//cfg.Transport = &http.Transport{
	//	Proxy: http.ProxyFromEnvironment,
	//	DialContext: (&net.Dialer{
	//		Timeout:   time.Duration(200) * time.Millisecond, //连接超时
	//		KeepAlive: 30 * time.Second,                      //长连接超时时间
	//	}).DialContext,
	//	ForceAttemptHTTP2:     true,                           //要使用自定义拨号程序或TLS配置并仍尝试HTTP/2升级，请将此设置为true。
	//	MaxIdleConns:          5,                              //最大空闲连接
	//	IdleConnTimeout:       time.Duration(3) * time.Second, //空闲超时时间
	//	TLSHandshakeTimeout:   1 * time.Second,                //tls握手超时时间
	//	ResponseHeaderTimeout: time.Duration(5) * time.Second, //如果非零，则指定在完全写入请求（包括其正文，如果有）之后等待服务器响应头的最长时间。 此时间不包括读响应体的时间。
	//}

	// 在查询服务列表时，如果没有找到任何可用的服务实例，就会等待 waitSecond 秒，直到有新的服务实例注册成功
	// 如果为0或者小于0，如果没找到，会立即返回错误
	cfg.WaitTime = config.WaitTime * time.Second

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return &ConsulServiceRegistryImpl{
		client: client,
		config: config,
	}
}
