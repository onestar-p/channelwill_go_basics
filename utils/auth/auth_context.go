/*
 * @Descripttion: 身份验证上下文相关
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-09 10:17:15
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-09 10:54:17
 */
package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthContext struct {
}

func NewAuthContext() *AuthContext {
	return &AuthContext{}
}

type userIDkey struct{}

/**
 * @name: 用户ID写入上下文
 * @param {context.Context} c
 * @param {string} uid
 * @return {*}
 */
func (ac *AuthContext) ContextWithUserId(c context.Context, uid string) context.Context {
	return context.WithValue(c, userIDkey{}, uid)
}

/**
 * @name: 通过context上下文获取userid
 * @param {context.Context} c
 * @return {*}
 */
func (ac *AuthContext) UserIDFromContext(c context.Context) (string, error) {
	v := c.Value(userIDkey{})
	uid, ok := v.(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}
	return uid, nil
}
