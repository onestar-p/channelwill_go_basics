/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-08 16:03:23
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-09 10:30:17
 */
package interceptor

import (
	"channelwill_go_basics/utils"
	"channelwill_go_basics/utils/token"
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader = "authorization"
	bearerPrefix        = ""
)

type Auth struct {
}

func (a *Auth) Do(c context.Context) (context.Context, error) {
	tkn, err := a.tokenFromContext(c)
	if err != nil {
		return c, fmt.Errorf("cannot parse public key:%v", err)
	}

	// 验证加密后的数据，并解密得到结果
	if tkn != "" {
		uid, err := utils.NewToken(token.JWTType).Verify(tkn)
		if err != nil { // 此处token验证根据实际业务情况做修改。
			return c, status.Errorf(codes.Unauthenticated, "token not valid:%v", err)
		}
		// 写入上下文
		c = utils.Auth.ContextWithUserId(c, uid)
	}
	return c, nil
}

/**
 * @name: 通过上下文获取token~
 * @param {context.Context} c
 * @return {*}
 */
func (a *Auth) tokenFromContext(c context.Context) (string, error) {
	unauthenticated := status.Error(codes.Unauthenticated, "")

	// 获取请求中的header数据
	m, ok := metadata.FromIncomingContext(c)
	if !ok {
		return "", unauthenticated
	}
	tkn := ""

	// 截取用户验证token
	for _, v := range m[authorizationHeader] {
		if strings.HasPrefix(v, bearerPrefix) {
			tkn = v[len(bearerPrefix):]
		}
	}
	// if tkn == "" {
	// 	return "", unauthenticated // 登录失败
	// }

	return tkn, nil
}
