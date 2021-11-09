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

	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
}

func (s *Service) SayHello(c context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	fmt.Println("Hellow")
	return &emptypb.Empty{}, nil
}
