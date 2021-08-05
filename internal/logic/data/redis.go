package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/php403/im/internal/logic/biz"
	"time"
)

const (
	UserTokenCacheKey string = "user:token:"
	UidCometServiceIdCacheKey string = "uid:comet:serverId:"
)

const DefaultExpireTime = 2*time.Hour


func GetAppServerName() string {
	return "logic"
}

func GetCacheKey(key string) string {
	return GetAppServerName()+":"+key+":"
}

func(ur * UserRepo) GetTokenCache(ctx context.Context,uniqueId string) (token string,err error) {
	token,err = ur.data.rdb.Get(ctx,GetCacheKey(uniqueId)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return
}

func(ur * UserRepo) SetTokenCache(ctx context.Context,uniqueId string,token string) (err error) {
	_,err = ur.data.rdb.Set(ctx,GetCacheKey(uniqueId),token,DefaultExpireTime).Result()
	return
}

func(ur * UserRepo) SetUserInfoCache(ctx context.Context,token string,user *biz.User) (err error)  {
	_,err= ur.data.rdb.HSet(ctx,token, map[string]interface{}{"id": user.Id, "app_id": user.AppId, "area_id":user.AreaId,"team_id":user.TeamId}).Result()
	return
}

func(ur *UserRepo) GetUserInfoCache(ctx context.Context,token string) (res map[string]string,err error) {
	res,err = ur.data.rdb.HGetAll(ctx,GetCacheKey(token)).Result()
	if err == redis.Nil {
		return res, nil
	}
	return
}




