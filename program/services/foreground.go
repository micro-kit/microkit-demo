package services

import (
	"context"
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/micro-kit/micro-common/cache"
	"github.com/micro-kit/micro-common/crypto"
	"github.com/micro-kit/micro-common/microerror"
	"github.com/micro-kit/microkit-demo/internal/pb"
	"github.com/micro-kit/microkit-demo/program/models"
)

/* 提供给客户端使用的rpc */

// Foreground 实现grpc客户端rpc接口
type Foreground struct {
	Base
}

// Login 登录接口
func (f *Foreground) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	username := req.Username
	password := req.Password
	user := new(models.User)
	password = crypto.PasswordHash(password)
	err := user.UserLogin(username, password)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = microerror.GetMicroError(10001)
		}
		return nil, err
	}
	// 将登录信息缓存入redis
	key, token := cache.GetUserLoginToken()
	userJs, _ := json.Marshal(user)
	err = cache.GetClient().Set(key, string(userJs), 0).Err()
	if err != nil {
		err = microerror.GetMicroError(10000, err)
		return nil, err
	}

	user.Token = token
	user.Password = ""

	reply := &pb.LoginReply{
		Username: user.Username,
		Nickname: user.Nickname,
		Password: user.Password,
		Token:    user.Token,
	}

	return reply, nil
}
