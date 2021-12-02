package utils_test

import (
	"channelwill_go_basics/global"
	"channelwill_go_basics/initialize"
	"channelwill_go_basics/utils"
	"channelwill_go_basics/utils/encrypter"
	"channelwill_go_basics/utils/id_generate"
	"channelwill_go_basics/utils/token"
	"fmt"
	"testing"
)

// ID生成测试
func TestIDGenerate(t *testing.T) {
	initialize.InitConfig()

	// 雪花ID
	id, _ := utils.IDGenerate(id_generate.Snowflake).GetID()
	fmt.Printf("Snowflake ID: %s\n", id)

	// uuid
	id, _ = id_generate.NewIDGenerate(id_generate.GenUuid).GetID()
	fmt.Printf("UUID: %s\n", id)

}

// 生成token测试
func TestGentToken(t *testing.T) {
	token, _ := utils.NewToken(token.SHA1Type).GenToken("my.shopify.com", 0)
	fmt.Printf("Token:%v\n", token)

}

//AES 测试
func TestEncrypter(t *testing.T) {

	res, _ := utils.AES.Encrypt("123")
	fmt.Printf("AES Encrypt:%v\n", res)

	res, _ = utils.AES.Decrypt(res)
	fmt.Printf("AES Decrypt:%v\n", res)

	fmt.Println("======")
	aes := encrypter.NewAES(&encrypter.AesConfig{
		Key:  global.ApplicationConfig.AESInfo.Key,  // 秘钥
		Iv:   global.ApplicationConfig.AESInfo.Iv,   // 秘钥偏移量
		Mode: global.ApplicationConfig.AESInfo.Mode, // 数据格式
	})
	res, _ = aes.Encrypt("abc")
	fmt.Printf("AES Encrypt:%v\n", res)

	res, _ = aes.Decrypt(res)
	fmt.Printf("AES Decrypt:%v\n", res)

}
