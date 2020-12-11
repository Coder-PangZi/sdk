package dingtalkrobot

import (
	"encoding/json"
	"io"
	"unsafe"
)

// 基础消息体
type Message struct {
	Typ      string          `json:"msgtype"`
	content  MessageContent  // 类型消息体
	AT       *at             `json:"at,omitempty"`
	Text     TextContent     `json:"text,omitempty"`
	Link     LinkContent     `json:"link,omitempty"`
	Markdown MarkdownContent `json:"markdown,omitempty"`
}

type at struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

func NewMessage() *Message {
	return &Message{
		AT: &at{AtMobiles: make([]string, 0)},
	}
}

func (m *Message) Bytes() []byte {
	m.content.Apply(m)
	buf, _ := json.Marshal(m)
	return buf
}

func (m *Message) String() string {
	buf := m.Bytes()
	return *(*string)(unsafe.Pointer(&buf))
}

func (m *Message) AddAtMobiles(mobiles ...string) {
	m.AT.AtMobiles = append(m.AT.AtMobiles, mobiles...)
}

func (m *Message) SetIsAtAll(isAtAll bool) {
	m.AT.IsAtAll = isAtAll
}

func (m *Message) SetContent(content MessageContent) {
	m.content = content
}

// 响应
type Response struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}

func DecodeResponse(code io.ReadCloser) (*Response, error) {
	rsp := new(Response)
	err := json.NewDecoder(code).Decode(rsp)
	if err != nil {
		return nil, err
	}
	return rsp, err
}

// 消息类型接口
type MessageContent interface {
	Type() string
	Apply(*Message)
}

// 纯文本消息
type TextContent struct {
	Content string `json:"content"`
}

func (t TextContent) Type() string {
	return "text"
}

func NewTextContent(content string) MessageContent {
	t := new(TextContent)
	t.Content = content
	return t
}

func NewTextMessage(content string) *Message {
	msg := NewMessage()
	msg.SetContent(NewTextContent(content))
	return msg
}

func (t TextContent) Apply(m *Message) {
	m.Text = t
	m.Typ = t.Type()
}

// 链接类消息

type LinkContent struct {
	Text       string `json:"text,omitempty"`
	Title      string `json:"title,omitempty"`
	PicUrl     string `json:"picUrl,omitempty"`
	MessageUrl string `json:"messageUrl,omitempty"`
}

func (l LinkContent) Type() string {
	return "link"
}

func (l LinkContent) Apply(msg *Message) {
	msg.Typ = l.Type()
	msg.Link = l
	msg.AT = nil
}

func NewLinkContent(text, title, picUrl, messageUrl string) MessageContent {
	return &LinkContent{Text: text, Title: title, PicUrl: picUrl, MessageUrl: messageUrl}
}

func NewLinkMessage(text, title, picUrl, messageUrl string) *Message {
	msg := NewMessage()
	msg.SetContent(NewLinkContent(text, title, picUrl, messageUrl))
	return msg
}

// Markdown 消息
type MarkdownContent struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

func (m MarkdownContent) Type() string {
	return "markdown"
}

func (m MarkdownContent) Apply(msg *Message) {
	msg.Typ = m.Type()
	msg.Markdown = m
	msg.AT = nil
}

func NewMarkdownContent(text, title string) MessageContent {
	return &MarkdownContent{Text: text, Title: title}
}

func NewMarkdownMessage(text, title string) *Message {
	msg := NewMessage()
	msg.SetContent(NewMarkdownContent(text, title))
	return msg
}
