package interceptor

import (
	"channelwill_go_basics/utils/jwt"
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
	PublicKeyFile string
}

func (a *Auth) Do(c context.Context) (context.Context, error) {
	tkn, err := a.tokenFromContext(c)
	if err != nil {
		return c, fmt.Errorf("cannot parse public key:%v", err)
	}

	// 验证加密后的数据，并解密得到结果
	if tkn != "" {
		// 解析公钥
		pubKey, err := jwt.NewJWTKey(a.PublicKeyFile).GetPublicKey()
		if err != nil {
			return c, fmt.Errorf("cannot parse public key:%v", err)
		}
		verifier := &jwt.JWTTokenVerifyer{
			PublicKey: pubKey,
		}
		uid, err := verifier.Verify(tkn)
		if err != nil {
			return c, status.Errorf(codes.Unauthenticated, "token not valid:%v", err)
		}
		// 写入上下文
		c = ContextWithUserId(c, uid)
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

type userIDkey struct{}

/**
 * @name: 用户ID写入上下文
 * @param {context.Context} c
 * @param {string} uid
 * @return {*}
 */
func ContextWithUserId(c context.Context, uid string) context.Context {
	return context.WithValue(c, userIDkey{}, uid)
}

/**
 * @name: 通过context上下文获取userid
 * @param {context.Context} c
 * @return {*}
 */
func UserIDFromContext(c context.Context) (string, error) {
	v := c.Value(userIDkey{})
	uid, ok := v.(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "")
	}
	return uid, nil
}
