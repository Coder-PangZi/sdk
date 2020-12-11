# 钉钉自定义机器人

[钉钉文档](https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq)

## 使用

[参考](./robot_test.go)

```go

rob := NewCustomRobot(
    WithSecureType(dingtalkrobot.KeyWord),
    WithAccessToken("xxxxx"))

rsp, err := rob.Send(NewTextMessage("【GM】hahaha"))
if err != nil {
    t.Fatal(err)
}
fmt.Printf("%#v", rsp)
```