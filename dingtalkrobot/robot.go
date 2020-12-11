package dingtalkrobot

import (
	"bytes"
	"errors"
	"net/http"
	"time"
)


type Robot struct {
	cfg *Config
	cli *http.Client
}

func NewCustomRobot(opts ...Option) *Robot {

	cli := &http.Client{
		Transport: &http.Transport{
			//TLSHandshakeTimeout:   0,
			//IdleConnTimeout:       1 * time.Second,
			//ResponseHeaderTimeout: 0,
			//ExpectContinueTimeout: 0,
		},
		Timeout: 5 * time.Second,
	}
	return &Robot{cfg: NewConfig(opts...), cli: cli}
}

func (rob *Robot) Send(msg *Message) (*Response, error) {
	ddRsp, err := rob.cli.Post(rob.cfg.Addr(), ContentType, bytes.NewReader(msg.Bytes()))
	if ddRsp != nil {
		defer ddRsp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if ddRsp.StatusCode != 200 {
		return nil, errors.New("status code not 200")
	}
	rsp, err := DecodeResponse(ddRsp.Body)

	if err != nil {
		return nil, err
	}
	return rsp, nil
}
