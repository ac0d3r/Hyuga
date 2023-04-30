package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// LarkData generate sign for lark
func LarkData(subject, content, secret string) ([]byte, error) {
	var (
		sign string
		err  error
	)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	if secret != "" {
		sign, err = larkSign(timestamp, secret)
		if err != nil {
			return nil, err
		}
	}
	data, err := buildLarkPostData(subject, content, timestamp, sign)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func larkSign(timestamp, secret string) (string, error) {
	key := fmt.Sprintf("%s\n%s", timestamp, secret)
	var data []byte
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

func buildLarkPostData(subject, content, timestamp, sign string) ([]byte, error) {
	msg := contentMsg{
		Tag:  "text",
		Text: content,
	}
	pd := &postData{
		Timestamp: timestamp,
		Sign:      sign,
		MsgType:   "post",
	}
	pd.Content.Post.ZhCN.Title = subject
	pd.Content.Post.ZhCN.Content = append(pd.Content.Post.ZhCN.Content, []contentMsg{msg})
	data, err := json.Marshal(pd)
	if err != nil {
		return nil, err
	}
	return data, err
}

type postData struct {
	Timestamp string  `json:"timestamp"`
	Sign      string  `json:"sign,omitempty"`
	MsgType   string  `json:"msg_type"`
	Content   content `json:"content"`
}

type content struct {
	Post post `json:"post"`
}

type zhCN struct {
	Title   string         `json:"title"`
	Content [][]contentMsg `json:"content"`
}

type post struct {
	ZhCN zhCN `json:"zh_cn,omitempty"`
}

type contentMsg struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}
