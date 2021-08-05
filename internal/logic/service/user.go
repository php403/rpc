package service

import (
	"context"
	"github.com/gin-gonic/gin"
	pb "github.com/php403/im/api/logic"
	"github.com/php403/im/internal/logic/app"
	"github.com/php403/im/internal/logic/biz"
	"github.com/php403/im/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserRequest struct {
	AppId     int64  `form:"app_id" binding:"required"`
	AreaId    int64  `form:"area_id" binding:"required"`
	UniqueKey string `form:"unique_key" binding:"required"`
}

type UserReply struct {
	Token	string `json:"token"`
}

func NewUserService(user *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		user: user,
		log: logger,
	}
}

func (s *UserService) Auth(c *gin.Context)  {
	param := &UserRequest{}
	err := c.ShouldBind(param)
	if err != nil {
		s.log.Log(log.LevelError,"logic","auth bind param err!")
		app.Errors(c,log.InvalidParams)
		return
	}
	token,err := s.user.Auth(context.Background(), param.UniqueKey)
	if err != nil {
		s.log.Log(log.LevelError,"logic","token err!")
		app.Errors(c,log.ServerError)
	}
	app.Result(c,map[string]string{"token":token},app.OK)
	return
}



/*//todo server name
func (s *UserService) Connect(ctx context.Context, req *pb.ConnectRequest) (*pb.ConnectReply, error) {
	appid, areaid, room,uid ,err := s.user.GetTeamInfo(ctx,req.Server,req.Token)
	if err != nil {
		return &pb.ConnectReply{}, err
	}
	return &pb.ConnectReply{Server: "test",Appid: appid,AreaId: areaid, RoomId: room, Uid: uid}, nil
}

func (s *UserService) DisConnect(ctx context.Context, req *pb.DisconnectRequest) (*pb.DisconnectReply, error) {
	return &pb.DisconnectReply{},nil
}*/

func(s *UserService) Receive (ctx context.Context, req *pb.ReceiveRequest) (receiveReply *pb.ReceiveReply, err error) {
	//todo metadata get uid
	md, ok := metadata.FromIncomingContext(ctx)
	if!ok{
		return receiveReply, status.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
	val,ok := md["token"]
	if !ok {
		return receiveReply, status.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
	err = s.user.ReceiveMsg(ctx,req.Server,val[0],req.Proto)
	return &pb.ReceiveReply{},nil
}
