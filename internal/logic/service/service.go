package service

import (

	"github.com/google/wire"
	"github.com/php403/im/api/logic"
	"github.com/php403/im/internal/logic/biz"
	"github.com/php403/im/pkg/log"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserService)

type Reply struct {
	code uint64
	message string
}

type UserService struct {
	logic.UnimplementedLogicServer

	log log.Logger

	user *biz.UserUsecase
}