package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ac0d3r/hyuga/thirdparty/notify/crypto"
)

type DingTalk struct {
}

func NewDingTalk() *DingTalk {
	return &DingTalk{}
}

func (p *DingTalk) Send(token, secret, subject, content string) error {
	if len(token) == 0 || len(secret) == 0 {
		return fmt.Errorf("dingtalk token or secret is empty")
	}

	data, err := buildPostData(subject, content)
	if err != nil {
		return err
	}
	url, err := crypto.DingTalkURL(token, secret)
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

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("statusCode: %d, body: %v", resp.StatusCode, string(result))
		return fmt.Errorf("dingtalk message response error %w", err)
	}
	return nil
}

func buildPostData(subject, content string) ([]byte, error) {
	content = fmt.Sprintf("### %s\n>%s", subject, content)
	type postData struct {
		MsgType  string `json:"msgtype"`
		Markdown struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		} `json:"markdown"`
	}
	pd := &postData{MsgType: "markdown"}
	pd.Markdown.Title = subject
	pd.Markdown.Text = content
	data, err := json.Marshal(pd)
	if err != nil {
		return nil, err
	}
	return data, err
}
