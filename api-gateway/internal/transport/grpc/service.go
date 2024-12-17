package grpc

import "gitlab.crja72.ru/gospec/go9/netevent/api-gateway/pkg/api/gateway"

type GatewayServer struct {
	gateway.UnimplementedGatewayServer
}

func NewGatewayServer() *GatewayServer {
	return &GatewayServer{}
}
