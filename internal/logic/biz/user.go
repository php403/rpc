package biz

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/php403/im/api/logic"
	"github.com/php403/im/api/protocol"
	"github.com/php403/im/pkg/log"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type UserStatus int32

const (
	RESERVE UserStatus = iota
	NORMAL
	FROZEN
)

type User struct {
	Id				int64			`gorm:"primaryKey" json:"id"`
	UniqueKey		string			`gorm:"unique" json:"unique_key"`
	AppId			int64			`json:"app_id"`
	AreaId 			int64			`json:"area_id"`
	TeamId 			int64			`json:"team_id"`
	Status			UserStatus		`gorm:"default:1" json:"status"`
	CreateTime 		uint32			`gorm:"autoCreateTime" json:"create_time"`
	LastLoginTime 	uint32			`gorm:"autoCreateTime" json:"last_login_time"`
	Like      		int64			`gorm:"-"`
}


type UserRepo interface {
	GetUser(ctx context.Context,UniqueKey string) (*User ,error)
	CreateUser(ctx context.Context,user *User) error

	GetTokenCache(ctx context.Context,UniqueKey string) (string ,error)
	SetTokenCache(ctx context.Context,UniqueKey string,Token string) error

	SetUserInfoCache(ctx context.Context,Token string,user *User) error
	GetUserInfoCache(ctx context.Context,token string) (res map[string]string,err error)

	SendMsg(topic string,body []byte) (err error)
}

type UserUsecase struct {
	repo UserRepo
	log log.Logger
}

type UserToken struct {
	Uid 	int64			`json:"uid"`
	AppId 	int64			`json:"app_id"`
	AreaId 	int64			`json:"area_id"`
	TeamId 	int64			`json:"team_id"`
	ExpireTime time.Time	`json:"expire_time"`

}

func NewUsercase(repo UserRepo,logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo,log: logger}
}

func (uc *UserUsecase) Auth(ctx context.Context,UniqueKey string)  (token string ,err error) {
	var user *User
	if user,err = uc.repo.GetUser(ctx,UniqueKey); err == nil{
		//生成token
		var byteToken []byte
		byteToken,err = CreateToken(user)
		if err != nil {
			return
		}
		_ = uc.repo.SetTokenCache(ctx,UniqueKey,string(byteToken))
		token = string(byteToken)
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	user = &User{
		UniqueKey : UniqueKey,
		AppId:1,
		AreaId: 1,
		Status:NORMAL,
	}
	err = uc.repo.CreateUser(ctx,user)
	if err != nil {
		return
	}
	byteToken,err := CreateToken(user)
	if err != nil {
		return
	}
	token = string(byteToken)
	_ = uc.repo.SetTokenCache(ctx,UniqueKey,token)
	return
}

func CreateToken(user *User) (byteToken []byte,err error) {
	expireTime,_ := time.ParseDuration("10h")
	userToken := &UserToken{Uid: user.Id,AppId: user.AppId,AreaId: user.AreaId,TeamId: user.TeamId,ExpireTime:time.Now().Add(expireTime)}
	byteToken,err = json.Marshal(userToken)
	return
}

func (uc *UserUsecase) GetTeamInfo(ctx context.Context, server string, token string)(appid int64,
	areaid int64, room int64,uid int64,err error) {
	userInfo,err := uc.repo.GetUserInfoCache(ctx,token)
	if err != nil {
		return
	}
	if len(userInfo) == 0 {

		return
	}
	if appid,err = strconv.ParseInt(userInfo["app_id"], 10, 64); err != nil{
		return
	}
	if areaid,err = strconv.ParseInt(userInfo["area_id"], 10, 64);err != nil{
		return
	}
	if room,err = strconv.ParseInt(userInfo["team_id"], 10, 64); err != nil{
		return
	}
	if uid,err =  strconv.ParseInt(userInfo["uid"], 10, 64); err != nil {
		return
	}
	return
}

func(uc *UserUsecase) ReceiveMsg(ctx context.Context,server string,token string,protoc *protocol.Proto) (err error) {
	//todo 系统化topicname
	var userToken UserToken
	err = json.Unmarshal([]byte(token),&userToken)
	pb := &logic.PushMsg{
		Msg: protoc.Body,
		AppId: strconv.FormatInt(userToken.AppId,10),
		RoomId: strconv.FormatInt(userToken.AreaId,10),
		UserId: userToken.Uid,
		Server:server,
		Operation:protoc.Op,
	}
	switch protoc.Op {
	case protocol.OpSendMsg:
		pb.Type = logic.PushMsg_PUSH
	case protocol.OpSendAppMsg:
		pb.Type = logic.PushMsg_BROADCAST
	case protocol.OpSendRoomMsg:
		pb.Type = logic.PushMsg_ROOM
	default:
		return errors.New("msg op err!")
	}
	msg,errMarshal := proto.Marshal(pb)
	if errMarshal != nil {
		return err
	}
	err = uc.repo.SendMsg("gameIm"+server,msg)
	return
}



