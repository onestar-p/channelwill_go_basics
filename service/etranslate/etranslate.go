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
	"channelwill_go_basics/utils/auth"
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
}

func (s *Service) SayHello(c context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	uid, _ := auth.NewAuthContext().UserIDFromContext(c)
	fmt.Println("Hellow", uid)
	return &emptypb.Empty{}, nil
}
