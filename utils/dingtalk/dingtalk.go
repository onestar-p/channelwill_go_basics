package dingtalk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"channelwill_go_basics/global"
)

type DingtalkInterfaced interface {
	SetAt([]string)
	Send() error
}

const url = "https://oapi.dingtalk.com/robot/send?access_token="

type Text struct {
	Content string `json:"content"`
}
type At struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

// text 消息结构体
type textMsg struct {
	Text    Text   `json:"text"`
	MsgType string `json:"msgtype"`
	At      At     `json:"at"`
}

func NewTextMessage(content string) *textMsg {
	msg := &textMsg{
		Text:    Text{Content: content},
		MsgType: "text",
	}
	return msg
}
func (t textMsg) Marshal() []byte {
	j, _ := json.Marshal(t)
	return j
}

// 钉钉机器人
type Dingtalk struct {
	AccessToken string
}

func NewDingtalk() *Dingtalk {
	return &Dingtalk{
		AccessToken: global.ApplicationConfig.DingtalkInfo.AccessToken,
	}
}

type messageFunc func() []byte

func (d Dingtalk) SendMessage(msgFunc messageFunc) error {
	msg := msgFunc() // 获取消息内容

	// 发起请求
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", url, d.AccessToken), strings.NewReader(string(msg)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解析返回结果
	type result struct {
		Errcode int
		Errmsg  string
	}
	r := result{}
	json.Unmarshal(body, &r)
	if r.Errcode > 0 {
		return fmt.Errorf(r.Errmsg)
	}
	return nil
}
