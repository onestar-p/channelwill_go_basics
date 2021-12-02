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
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"channelwill_go_basics/forms"
	authpb "channelwill_go_basics/proto/auth/gen/v1"
	"channelwill_go_basics/utils"
	"channelwill_go_basics/utils/token"
)

type Service struct {
	TokenGenerator TokenGenerator
	TokenExpire    time.Duration
}

type TokenGenerator interface {
	GenToken(userId string, expire time.Duration) (string, error)
}

func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {

	login := forms.AuthLoginForm{
		UserName: req.UserName,
		Passwd:   req.Passwd,
	}
	if err := utils.Validate.Verify(login); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	aid := "123123123"

	fmt.Println("Auth Login...")
	sha1Token, _ := utils.NewToken(token.SHA1Type).GenToken("my.shopify.com", 0)
	fmt.Println(sha1Token)

	// 生成token
	tkn, err := utils.NewToken(token.JWTType).GenToken(aid, s.TokenExpire)
	if err != nil {
		zap.S().Error("cannot generate token", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	return &authpb.LoginResponse{
		UserToken: tkn,
	}, nil
}
