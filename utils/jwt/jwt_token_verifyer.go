package jwt

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type JWTTokenVerifyer struct {
	PublicKey *rsa.PublicKey
}

/**
 * @name: 验证token，并返回内容
 * @param {string} token
 * @return {string, error}
 */
func (v *JWTTokenVerifyer) Verify(token string) (string, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return v.PublicKey, nil
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
