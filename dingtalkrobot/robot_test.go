package dingtalkrobot

import (
	"io/ioutil"
	"testing"
	"unsafe"
)

func bytes2str(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}


// 关键词验证
// 需要在消息体中加入预设的关键词
func TestRobot_Send(t *testing.T) {

	buf, err := ioutil.ReadFile("./accesstoken.hide")
	if err != nil {
		t.Fatal("access token file not exist",err)
	}

	rob := NewCustomRobot(
		WithSecureType(KeyWord),
		WithAccessToken(bytes2str(buf)))

	msg := NewTextMessage("【GM】hahaha")
	t.Log("msg", msg.String())
	rsp, err := rob.Send(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", rsp)
}

// 签名验证方式
func TestRobot_Send2(t *testing.T) {

	buf, err := ioutil.ReadFile("./accesstoken_sign.hide")
	if err != nil {
		t.Fatal("access token file not exist",err)
	}

	buf2, err := ioutil.ReadFile("./secret.hide")
	if err != nil {
		t.Fatal("secret file not exist",err)
	}

	rob := NewCustomRobot(
		WithSecureType(Sign),
		WithSecret(bytes2str(buf2)),
		WithAccessToken(bytes2str(buf)))
	msg := NewMessage()
	content := TextContent{Content: "hahaha"}
	msg.SetContent(content)
	t.Log("msg", msg.String())
	rsp, err := rob.Send(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", rsp)
}
