package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Bark struct {
	url string
}

const (
	defaultBarkServer = "api.day.app"
)

func NewBark() *Bark {
	return &Bark{
		url: fmt.Sprintf("https://%s/push", defaultBarkServer),
	}
}

type barkPostData struct {
	DeviceKey string `json:"device_key"`
	Title     string `json:"title"`
	Body      string `json:"body,omitempty"`
	Badge     int    `json:"badge,omitempty"`
	Sound     string `json:"sound,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Group     string `json:"group,omitempty"`
	URL       string `json:"url,omitempty"`
}

func (p *Bark) Send(key, server, subject, content string) error {
	if len(key) == 0 {
		return fmt.Errorf("bark key is empty")
	}
	if server != "" {
		p.url = fmt.Sprintf("https://%s/push", server)
	}

	pd := &barkPostData{
		DeviceKey: key,
		Title:     subject,
		Body:      content,
		Sound:     "alarm.caf",
	}
	data, err := json.Marshal(pd)
	if err != nil {
		return err
	}

	resp, err := http.Post(p.url, "application/json; charset=utf-8", bytes.NewReader(data))
	if err != nil {
		return err
	}
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("send bark message failed %w",
			fmt.Errorf("statusCode: %d, body: %v", resp.StatusCode, string(result)))
	}
	return nil
}
