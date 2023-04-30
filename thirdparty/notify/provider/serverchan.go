package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultServer = "sctapi.ftqq.com"

type ServerChan struct {
}

func NewServerChan() *ServerChan {
	return &ServerChan{}
}

func (p *ServerChan) Send(userID, sendKey, subject, content string) error {
	if len(userID) == 0 || len(sendKey) == 0 {
		return fmt.Errorf("server chan userID or sendKey is empty")
	}

	url := fmt.Sprintf("https://%s/%s.send", defaultServer, sendKey)
	type postData struct {
		Text string `json:"text"`
		Desp string `json:"desp"`
	}
	pd := &postData{
		Text: subject,
		Desp: content,
	}

	data, err := json.Marshal(pd)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("send server chan failed %w", fmt.Errorf("http status code: %d", resp.StatusCode))
	}
	return nil
}
