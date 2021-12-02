package token

import (
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"

	etRoot "channelwill_go_basics"
	"channelwill_go_basics/global"
	"channelwill_go_basics/utils/id_generate"
	utilJwt "channelwill_go_basics/utils/jwt"
)

type TokenInterfaced interface {
	GenToken(str string, expire time.Duration) (string, error)
	Verify(token string) (string, error)
}

type TokenType int

const (
	JWTType  TokenType = iota // JWT token生成
	SHA1Type                  // 使用SHA1加密算法生成 token
)

var (
	jwtPrivateKey *rsa.PrivateKey
	jwtPpublicKey *rsa.PublicKey
)

func NewToken(tokenType TokenType) TokenInterfaced {
	// JWT Token
	if tokenType == JWTType {
		appConfig := global.ApplicationConfig
		if jwtPrivateKey == nil {
			privateKeyFile := etRoot.Path("config/cert/private.key") // 私钥路径
			privKey, err := utilJwt.NewJWTKey(privateKeyFile).GetPrivateKey()
			if err != nil {
				zap.S().Fatal("cannot pare private key", zap.Error(err))
			}
			jwtPrivateKey = privKey
		}

		if jwtPpublicKey == nil {
			publicFile := etRoot.Path("config/cert/public.key") // 公钥路径
			publKey, err := utilJwt.NewJWTKey(publicFile).GetPublicKey()
			if err != nil {
				zap.S().Fatal("cannot pare public key", zap.Error(err))
			}
			jwtPpublicKey = publKey
		}

		return NewTokenJWT(JWTConfig{
			Issuer:     appConfig.JwtInfo.Issuer,
			PrivateKey: jwtPrivateKey,
			PublicKey:  jwtPpublicKey,
		})
	}

	// SHA1
	if tokenType == SHA1Type {
		return SHA1Token{}
	}

	return nil

}

// Generate Token by SHA1
type SHA1Token struct{}

func (st SHA1Token) GenToken(str string, expire time.Duration) (string, error) {
	id, _ := id_generate.NewIDGenerate(id_generate.GenUuid).GetID()
	// MD5
	m5 := md5.New()
	m5.Write([]byte(id))
	md5Str := hex.EncodeToString(m5.Sum(nil))
	// SHA1
	hStr := fmt.Sprintf("%s%s", md5Str, str)
	h := sha1.New()
	h.Write([]byte(hStr))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs), nil
}

func (st SHA1Token) Verify(token string) (string, error) {
	return "", nil
}

// Generate Token by JWT
type TokenJWT struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	isser      string
	nowFunc    func() time.Time
}

type JWTConfig struct {
	Issuer     string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewTokenJWT(config JWTConfig) *TokenJWT {
	return &TokenJWT{
		isser:      config.Issuer,
		privateKey: config.PrivateKey,
		publicKey:  config.PublicKey,
		nowFunc:    time.Now,
	}
}

/**
 * @name: 生成token
 * @msg:
 * @param {string} uid 用户ID，或者其他
 * @param {time.Duration} expire 有效时间
 * @return {token, error}
 */
func (j *TokenJWT) GenToken(uid string, expire time.Duration) (string, error) {
	nowSec := j.nowFunc().Unix()
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    j.isser,
		IssuedAt:  nowSec, // 什么时候颁发的
		ExpiresAt: nowSec + int64(expire.Seconds()),
		Subject:   uid,
	})

	// 签名
	return tkn.SignedString(j.privateKey)
}

/**
 * @name: 验证token，并返回内容
 * @param {string} token
 * @return {string, error}
 */
func (j *TokenJWT) Verify(token string) (string, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.publicKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("cannot parse token:%v", err)
	}
	if !t.Valid {
		return "", fmt.Errorf("token not valid")
	}

	clm, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("token claim is not StandardClaims")
	}

	if err := clm.Valid(); err != nil {
		return "", fmt.Errorf("claim not valid:%v", err)
	}

	return clm.Subject, nil
}
