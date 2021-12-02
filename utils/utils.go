package utils

import (
	"channelwill_go_basics/utils/auth"
	"channelwill_go_basics/utils/encrypter"
	"channelwill_go_basics/utils/id_generate"
	"channelwill_go_basics/utils/jwt"
	"channelwill_go_basics/utils/token"
	"channelwill_go_basics/utils/validate"
)

const (
	aesKey  = "k0b4o2t0l3t4a4n0" // 秘钥
	aesIv   = "d5o9m6a2d1p3u0l3" // 偏移量
	aesMode = "RD"               // 数据格式
)

var (
	AES        = encrypter.NewAES(&encrypter.AesConfig{Key: aesKey, Iv: aesIv, Mode: aesMode}) // AES
	IDGenerate = id_generate.NewIDGenerate                                                     // ID生成
	Auth       = auth.NewAuth()
	NewToken   = token.NewToken         // Token
	Validate   = validate.NewValidate() // 数据参数验证
	JWTKey     = jwt.NewJWTKey          // JWT 私钥/秘钥操作
)
