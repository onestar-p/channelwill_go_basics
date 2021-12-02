/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-04 16:37:08
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-09 09:35:18
 */
package etranslate

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"channelwill_go_basics/utils"
)

type Service struct {
}

func (s *Service) SayHello(c context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	uid, err := utils.Auth.UserIDFromContext(c)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "用户未登录")
	}
	fmt.Println("Hellow", uid)
	return &emptypb.Empty{}, nil
}
