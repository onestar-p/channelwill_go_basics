package encrypter

import (
	"encoding/base64"

	"github.com/forgoer/openssl"
)

var (
	padding = map[string]string{
		"RD": openssl.PKCS7_PADDING,
		"ZP": openssl.ZEROS_PADDING,
	}
)

// CBC AES加密
type Aes struct {
	key  []byte // 秘钥
	iv   []byte // 秘钥偏移量
	mode string // 数据格式

}

type AesConfig struct {
	Key  string // 秘钥；秘钥长度决定加密方式：16位：aes-128-cbc；24位：aes-192-cbc；32位：aes-256-cbc
	Iv   string // 秘钥偏移量
	Mode string // 数据格式
}

func NewAES(config *AesConfig) *Aes {
	return &Aes{
		key:  []byte(config.Key),
		iv:   []byte(config.Iv),
		mode: config.Mode,
	}
}

func (a Aes) getPadding(mode string) string {
	pad, ok := padding[mode]
	if ok {
		return pad
	} else {
		return padding["RD"]
	}
}

func (a Aes) Encrypt(data string) (string, error) {
	if dst, err := openssl.AesCBCEncrypt([]byte(data), a.key, a.iv, a.getPadding(a.mode)); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(dst), nil
	}
}

func (a Aes) Decrypt(dataBase64 string) (string, error) {
	dataByte, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		return "", err
	}

	if dst, err := openssl.AesCBCDecrypt(dataByte, a.key, a.iv, a.getPadding(a.mode)); err != nil {
		return "", err
	} else {
		return string(dst), nil
	}
}
