package server

import (
	"github.com/php403/im/api/logic"
	"github.com/php403/im/internal/logic/conf"
	"github.com/php403/im/internal/logic/service"
	"github.com/php403/im/pkg/log"
	"google.golang.org/grpc"
	"net"
)

type GrpcServer struct {
	Grpc	*grpc.Server
	GrpcLis net.Listener
	User 	*service.UserService
	C 		*conf.Server
}

//todo 后续封装grpc server start 放main函数
func NewGRPCServer(c *conf.Server,logger log.Logger,user *service.UserService) *GrpcServer {
	s := grpc.NewServer()
	logic.RegisterLogicServer(s,user)
	return &GrpcServer{
		Grpc: s,
		User: user,
		C:c,
	}
}




