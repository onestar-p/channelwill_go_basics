/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-08 14:28:12
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-09 09:32:24
 */
package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

type MiddlewareInterfaced interface {
	Do(c context.Context) (context.Context, error) // 执行中间件，关于错误处理，如非必要，请勿返回错误处理
}

type Interceptor struct {
	Middleware []MiddlewareInterfaced
}

func NewInterceptor() *Interceptor {
	return &Interceptor{}
}

func (i *Interceptor) AddMiddleware(m MiddlewareInterfaced) {
	i.Middleware = append(i.Middleware, m)
}

func (i *Interceptor) UnaryServerInterceptor() (grpc.UnaryServerInterceptor, error) {
	return i.HandleReq, nil
}

func (i *Interceptor) HandleReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	for _, m := range i.Middleware {
		ctx, err = m.Do(ctx)
		if err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}
