/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-04 23:08:53
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-08 17:13:15
 */
package auth

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"channelwill_go_basics/interceptor"
	authpb "channelwill_go_basics/proto/gen/auth/v1"
)

type Service struct {
	authpb.UnimplementedAuthServiceServer
	TokenGenerator    TokenGenerator
	AuthPublicKeyFile string // 公钥文件地址
	TokenExpire       time.Duration
}

type TokenGenerator interface {
	GenerateToken(userId string, expire time.Duration) (string, error)
}

func (s *Service) Login(c context.Context, req *emptypb.Empty) (*authpb.LoginResponse, error) {
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

func (s *Service) GetUserToken(c context.Context, req *emptypb.Empty) (*authpb.GetUserTokenResponse, error) {

	uid, err := interceptor.UserIDFromContext(c)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "用户未授权")
	}

	return &authpb.GetUserTokenResponse{
		Token: fmt.Sprintf("%s:%d", uid, 02),
	}, nil
}
