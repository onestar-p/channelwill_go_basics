/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-04 23:08:53
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-09 10:31:19
 */
package auth

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"channelwill_go_basics/forms"
	authpb "channelwill_go_basics/proto/gen/auth/v1"
	"channelwill_go_basics/utils/validate"
)

type Service struct {
	TokenGenerator    TokenGenerator
	AuthPublicKeyFile string // 公钥文件地址
	TokenExpire       time.Duration
}

type TokenGenerator interface {
	GenerateToken(userId string, expire time.Duration) (string, error)
}

func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {

	login := forms.AuthLoginForm{
		UserName: req.UserName,
		Passwd:   req.Passwd,
	}
	if err := validate.NewValidate().Verify(login); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	aid := "123123123"
	// 生成token
	tkn, err := s.TokenGenerator.GenerateToken(aid, s.TokenExpire)
	if err != nil {
		zap.S().Error("cannot generate token", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	return &authpb.LoginResponse{
		UserToken: tkn,
	}, nil
}
