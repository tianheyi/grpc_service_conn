package grpc_srv_conn

import (
	"fmt"
	"grpc_service_conn/service_registry"
	"testing"
)

var (
	TestConsulServiceName = "user-srv"
)

func TestGrpcServiceConfig(t *testing.T) {
	registry, _ := service_registry.NewServiceRegistry("consul", &service_registry.ServiceRegistryConfig{Host: "127.0.0.1", Port: 8500})
	config := &GrpcServiceConfig{LoadBalancingPolicy: "round_robin"}
	srvManger, _ := NewSrvManagerImpl(registry, config)
	conn, err := srvManger.GetGrpcConn(TestConsulServiceName, 0)
	if err != nil {
		t.Error(err.Error())
	}
	defer conn.Close()
	// do something
	fmt.Println(conn.GetState().String())
}
