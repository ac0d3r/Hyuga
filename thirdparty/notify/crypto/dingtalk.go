package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"
)

const defaultDingTalkURL = `https://oapi.dingtalk.com/robot/send`

func DingTalkURL(token, secret string) (string, error) {
	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	u, _ := url.Parse(defaultDingTalkURL)
	value := url.Values{}
	value.Set("access_token", token)

	if secret == "" {
		u.RawQuery = value.Encode()
		return u.String(), nil
	}

	sign, err := dingTalkSign(timestamp, secret)
	if err != nil {
		u.RawQuery = value.Encode()
		return u.String(), err
	}

	value.Set("timestamp", timestamp)
	value.Set("sign", sign)
	u.RawQuery = value.Encode()
	return u.String(), nil
}

func dingTalkSign(timestamp, secret string) (string, error) {
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	if _, err := io.WriteString(h, stringToSign); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
