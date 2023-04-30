package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ac0d3r/hyuga/thirdparty/notify/crypto"
)

const (
	larkAPI   = "https://open.larksuite.com/open-apis/bot/v2/hook/%s"
	feishuAPI = "https://open.feishu.cn/open-apis/bot/v2/hook/%s"
)

type Lark struct {
}

func NewLark() *Lark {
	return &Lark{}
}

func (p *Lark) Send(token, secret, subject, content string) error {
	return p.send(fmt.Sprintf(larkAPI, token), token, secret, subject, content)
}

func (p *Lark) SendWithFeishu(token, secret, subject, content string) error {
	return p.send(fmt.Sprintf(feishuAPI, token), token, secret, subject, content)
}

func (p *Lark) send(url, token, secret, subject, content string) error {
	if len(token) == 0 || len(secret) == 0 {
		return fmt.Errorf("lark token or secret is empty")
	}

	data, err := crypto.LarkData(subject, content, secret)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(data))
	if err != nil {
		return err
	}
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type response struct {
		Code          int         `json:"code"`
		Data          struct{}    `json:"data"`
		Msg           string      `json:"msg"`
		Extra         interface{} `json:"Extra"`
		StatusCode    int         `json:"StatusCode"`
		StatusMessage string      `json:"StatusMessage"`
	}

	var r response
	if err = json.Unmarshal(result, &r); err != nil {
		return err
	}

	if r.StatusMessage != "success" {
		return fmt.Errorf("lark message response %w", fmt.Errorf("body: %v", string(result)))
	}
	return nil
}
