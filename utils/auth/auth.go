/*
 * @Author: your name
 * @Date: 2021-11-12 17:37:56
 * @LastEditTime: 2021-11-13 17:54:39
 * @LastEditors: Please set LastEditors
 * @Description: 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 * @FilePath: /channelwill_go_basics/utils/auth/auth.go
 */
package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth struct{}

func NewAuth() *Auth {
	return &Auth{}
}

/**
 * @name: 用户ID写入上下文
 * @param {context.Context} c
 * @param {string} uid
 * @return {*}
 */
func (a *Auth) ContextWithUserId(c context.Context, uid string) context.Context {
	return context.WithValue(c, userIDkey{}, uid)
}

type userIDkey struct{}

/**
 * @name: 通过context上下文获取userid
 * @param {context.Context} c
 * @return {*}
 */
func (a *Auth) UserIDFromContext(c context.Context) (string, error) {
	v := c.Value(userIDkey{})
	uid, ok := v.(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}
	return uid, nil
}
