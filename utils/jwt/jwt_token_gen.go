package jwt

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTTokenGen struct {
	privateKey *rsa.PrivateKey
	isser      string
	nowFunc    func() time.Time
}

func NewJWTTokenGen(issuer string, privateKey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		isser:      issuer,
		nowFunc:    time.Now,
		privateKey: privateKey,
	}
}

/**
 * @name: 生成token
 * @msg:
 * @param {string} uid 用户ID，或者其他
 * @param {time.Duration} expire 有效时间
 * @return {token, error}
 */
func (t *JWTTokenGen) GenerateToken(uid string, expire time.Duration) (string, error) {
	nowSec := t.nowFunc().Unix()
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    t.isser,
		IssuedAt:  nowSec, // 什么时候颁发的
		ExpiresAt: nowSec + int64(expire.Seconds()),
		Subject:   uid,
	})

	// 签名
	return tkn.SignedString(t.privateKey)
}
