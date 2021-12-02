package email

import (
	"fmt"
	"testing"

	"channelwill_go_basics/initialize"
)

func TestEmail(t *testing.T) {
	initialize.InitConfig()
	initialize.InitRedis()
	to := []string{
		// "2912313265@qq.com",
		// "antxiaoye@gmail.com",
	}
	err := NewEmail().Send(to, "title test", "<h1>Hello test 你好测试</h1>")
	// email.
	fmt.Println(err)

}
