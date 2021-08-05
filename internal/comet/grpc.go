package comet

import (
	"context"
	"github.com/php403/im/api/logic"
	"github.com/php403/im/api/protocol"
	"google.golang.org/grpc/metadata"
)

//todo 简化流程 保留代码 后续可能需要改回
/*func (s *Server) Connect(ctx context.Context,token string) (Server string,
	Appid  int64, AreaId int64, RoomId int64, Uid  int64, err error) {
	reply, err := s.rpcClient.Connect(ctx, &logic.ConnectRequest{
		Server: s.serverID,
		Token:  token,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	return reply.Server,reply.Appid,reply.AreaId,reply.RoomId,reply.Uid,err
}

func (s *Server) DisConnect(ctx context.Context,uid int64) (err error) {
	_, err = s.rpcClient.Disconnect(ctx, &logic.DisconnectRequest{
		Server: s.serverID,
		Uid:  uid,
	})
	return
}

func (s *Server) Heartbeat(ctx context.Context) (err error) {
	_, err = s.rpcClient.Heartbeat(ctx, &logic.HeartbeatRequest{
		Server: s.serverID,
	})
	return
}*/

func (s *Server) Receive(ctx context.Context,token string,p *protocol.Proto) (err error) {
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)
	_,err = s.rpcClient.Receive(ctx,&logic.ReceiveRequest{Server: s.serverName,Proto: p})
	return
}