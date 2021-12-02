package dingtalk

import (
	"channelwill_go_basics/initialize"
	"fmt"
	"testing"
)

func TestDingtalk(t *testing.T) {
	initialize.InitConfig()
	err := NewDingtalk().SendMessage(func() []byte {
		return NewTextMessage("test222").Marshal()
	})
	fmt.Println(err)
}
