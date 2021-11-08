package interceptor

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"channelwill_go_basics/utils/jwt"
)

const (
	authorizationHeader = "authorization"
	bearerPrefix        = ""
)

func Auth(publicKeyFile string) (grpc.UnaryServerInterceptor, error) {
	// 解析公钥
	pubKey, err := jwt.NewJWTKey(publicKeyFile).GetPublicKey()
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key:%v", err)
	}

	i := &authInterceptor{
		verifier: &jwt.JWTTokenVerifyer{
			PublicKey: pubKey,
		},
	}
	return i.HandleReq, nil
}

type tokenVerifier interface {
	Verify(token string) (string, error)
}

type authInterceptor struct {
	verifier tokenVerifier
}

func (i *authInterceptor) HandleReq(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	tkn, err := tokenFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}

	// 验证加密后的数据，并解密得到结果
	if tkn != "" {
		uid, err := i.verifier.Verify(tkn)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token not valid:%v", err)
		}
		// 写入上下文
		ctx = ContextWithUserId(ctx, uid)
	}

	return handler(ctx, req)
}

/**
 * @name: 通过上下文获取token
 * @param {context.Context} c
 * @return {*}
 */
func tokenFromContext(c context.Context) (string, error) {
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
