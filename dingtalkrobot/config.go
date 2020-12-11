package dingtalkrobot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

type SecureType int

const (
	KeyWord SecureType = iota
	Sign

	Custom RobotType = iota

	DingDingRobotAddr = `https://oapi.dingtalk.com/robot/send`
	ContentType = `application/json; charset=utf-8`
)


type Config struct {
	st     SecureType
	addr   string
	at     string
	secret string
	rt     RobotType
}

func (c Config) Secret() string {
	return c.secret
}

func NewConfig(opts ...Option) *Config {
	cfg := new(Config)
	for i := range opts {
		opts[i].Apply(cfg)
	}
	return cfg
}

func (c Config) AccessToken() string {
	return c.at
}

func (c Config) Addr() string {
	if c.addr == "" {
		WithAddr(DingDingRobotAddr).Apply(&c)
	}
	addr := c.addr + "?access_token=" + c.AccessToken()
	if c.SecureType() == Sign {
		ts := time.Now().UnixNano() / 1e6
		addr += "&timestamp=" + strconv.FormatInt(ts, 10)
		sign := genSign(ts, c.Secret())
		addr += "&sign=" + sign
	}
	return addr
}

func genSign(ts int64, key string) string {
	keygen := hmac.New(sha256.New, []byte(key))
	keygen.Write([]byte(fmt.Sprintf("%d\n%s", ts, key)))

	return base64.StdEncoding.EncodeToString(keygen.Sum(nil))
}

func (c Config) SecureType() SecureType {
	return c.st
}

type Option interface {
	Apply(*Config)
}

type WithSecureType SecureType

func (st WithSecureType) Apply(config *Config) {
	config.st = SecureType(st)
}

type WithAddr string

func (addr WithAddr) Apply(config *Config) {
	config.addr = string(addr)
}

type WithAccessToken string

func (at WithAccessToken) Apply(config *Config) {
	config.at = string(at)
}

type WithSecret string

func (s WithSecret) Apply(config *Config) {
	config.secret = string(s)
}


type RobotType int

func (s RobotType) Apply(config *Config) {
	config.rt = s
}