package data

import (
	"context"
	"github.com/php403/im/internal/logic/biz"
	"github.com/php403/im/pkg/log"
)

type UserRepo struct {
	data *Data
	log  log.Logger
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &UserRepo{
		data: data,
		log:  logger,
	}
}

func(ur * UserRepo) CreateUser(ctx context.Context, user *biz.User) error {
	result := ur.data.db.WithContext(ctx).Create(user)
	return result.Error
}

func(ur * UserRepo) GetUser(ctx context.Context, uniqueId string) (user *biz.User,err error) {
	err = ur.data.db.WithContext(ctx).Where("unique_key = ?",uniqueId).First(&user).Error
	return
}





