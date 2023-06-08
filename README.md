## 结构
service_registry：用于建立consul客户端连接\
grpc_srv_conn：调用service_registry方法建立consul客户端连接，获取对应服务地址，根据相关负载均衡策略选择一个符合条件的地址

## grpc服务配置proto文件地址
service_config.proto文件：https://github.com/grpc/grpc-proto/blob/master/grpc/service_config/service_config.proto

## gRPC resolver与balancer通信验证
掘金地址：https://juejin.cn/post/7241575425028259896


