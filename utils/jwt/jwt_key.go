/*
 * @Descripttion: JWT相关操作
 */
package jwt

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type JWTKeyInterfaced interface {
	GetPrivateKey() (*rsa.PrivateKey, error)
	GetPublicKey() (*rsa.PublicKey, error)
}

type JWTKey struct {
	keyFileName string // 公钥或私钥文件地址
	JWTKeyInterfaced
}

func NewJWTKey(keyFileName string) *JWTKey {
	return &JWTKey{
		keyFileName: keyFileName,
	}
}

/**
 * @name: 获取私钥
 */
func (jt *JWTKey) GetPrivateKey() (*rsa.PrivateKey, error) {
	pkFile, err := os.Open(jt.keyFileName)
	if err != nil {
		zap.S().Fatal("cannot open pricate key ", zap.Error(err))
	}

	pkBytes, err := ioutil.ReadAll(pkFile) // 读取文件内容
	if err != nil {
		zap.S().Fatal("cannot read private key", zap.Error(err))
	}
	return jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
}

/**
 * @name: 获取公钥
 */
func (jt *JWTKey) GetPublicKey() (*rsa.PublicKey, error) {
	f, err := os.Open(jt.keyFileName)
	if err != nil {
		return nil, fmt.Errorf("cannot open public key : %v", err)
	}

	// 读取公钥全部内容
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key:%v", err)
	}

	// 解析公钥
	return jwt.ParseRSAPublicKeyFromPEM(b)
}
