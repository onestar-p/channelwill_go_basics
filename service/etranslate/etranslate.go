/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-04 16:37:08
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-08 09:22:07
 */
package etranslate

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	etranslatepb "channelwill_go_basics/proto/gen/etranslate/v1"
)

type Service struct {
	etranslatepb.UnimplementedEtranslateServiceServer
}

func (s *Service) SayHello(c context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	fmt.Println("Hellow")
	return &emptypb.Empty{}, nil
}
