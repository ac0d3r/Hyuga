# notifier

![CI](https://github.com/moonD4rk/notifier/workflows/CI/badge.svg?branch=main)

notifier is a simple Go library to send notification to other applications.

## Feature

| Provider                                                     | Code                                                         |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| [DingTalk](https://www.dingtalk.com/en)                      | [provider/dingtalk](https://github.com/moonD4rk/notifier/tree/main/provider/dingtalk) |
| [Bark](https://apps.apple.com/us/app/bark-customed-notifications/id1403753865) | [provider/bark](https://github.com/moonD4rk/notifier/tree/main/provider/bark) |
| [Lark](https://www.larksuite.com/en_us/)                     | [provider/lark](https://github.com/moonD4rk/notifier/tree/main/provider/lark) |
| [Feishu](https://www.feishu.cn/)                             | [provider/feishu](https://github.com/moonD4rk/notifier/tree/main/provider/feishu) |
| [Server 酱](https://sct.ftqq.com/)                           | [provider/serverchan](https://github.com/moonD4rk/notifier/tree/main/provider/serverchan) |

## Install

`go get -u github.com/moond4rk/notifier`

## Usage



```go
package main

import (
	"os"

	"github.com/moond4rk/notifier"
)

func main() {
	var (
		dingtalkToken     = os.Getenv("dingtalk_token")
		dingtalkSecret    = os.Getenv("dingtalk_secret")
		barkKey           = os.Getenv("bark_key")
		barkServer        = notifier.DefaultBarkServer
		feishuToken       = os.Getenv("feishu_token")
		feishuSecret      = os.Getenv("feishu_secret")
		larkToken         = os.Getenv("lark_token")
		larkSecret        = os.Getenv("lark_secret")
		serverChanUserID  = "" // server chan's userID could be empty
		serverChanSendKey = os.Getenv("server_chan_send_key")
	)
	notifier := notifier.New(
		notifier.WithDingTalk(dingtalkToken, dingtalkSecret),
		notifier.WithBark(barkKey, barkServer),
		notifier.WithFeishu(feishuToken, feishuSecret),
		notifier.WithLark(larkToken, larkSecret),
		notifier.WithServerChan(serverChanUserID, serverChanSendKey),
	)

	var (
		subject = "this is subject"
		content = "this is content"
	)
	if err := notifier.Send(subject, content); err != nil {
		panic(err)
	}
}
```

<img src="https://raw.githubusercontent.com/moonD4rk/staticfiles/master/picture/notifier-screenshot.png" width="480" align="left"/>
