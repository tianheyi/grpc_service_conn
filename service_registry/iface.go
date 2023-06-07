package service_registry

// 服务节点
type ServiceNode struct {
	Address string
	Port    int
}

// 服务注册中心配置信息
type ServiceRegistryConfig struct {
	Host string
	Port int
}

// 服务注册中心接口，负责服务注册、注销、服务发现
type ServiceRegistry interface {
	// Register 注册服务
	//  - id: 注册的服务对象的id，在服务中保证对象唯一性的标识
	//	- address: 注册的服务对象的ip地址
	// 	- port: 注册的服务对象的端口
	// 	- name: 要注册到的服务名称
	//  - tags: 服务对象标签
	Register(id string, address string, port int, name string, tags []string) error
	// DeRegister 注销服务
	//  - serviceId: 要注销的服务对象的id
	DeRegister(serviceId string) error
	// GetServiceList 获取服务列表 handler为nil则返回符合的服务列表，否则返回handle处理后的一个
	//  - service: 要获取的服务名称
	//  - tags: 过滤条件，通过标签过滤
	//  - passingOnly: 为true时只返回通过健康检查的服务列表
	//  - handler: 自定义的处理函数，可实现负载均衡等功能，不为nil时返回长度为1的服务列表
	//  弃用，使用grpc组件实现
	GetServiceList(service string, tags []string, passingOnly bool, handler func([]*ServiceNode) *ServiceNode) ([]*ServiceNode, error)
	// GetConfig 获取注册中心配置信息
	GetConfig() *ServiceRegistryConfig
}
