package notifier

import "github.com/ac0d3r/hyuga/thirdparty/notify/provider"

type Notifier struct {
	Bark       *provider.Bark
	DingTalk   *provider.DingTalk
	Lark       *provider.Lark
	ServerChan *provider.ServerChan
}

func NewNotify() *Notifier {
	n := &Notifier{
		Bark:     provider.NewBark(),
		DingTalk: provider.NewDingTalk(),
	}
	return n
}

var n = NewNotify()

func WithBark(key, server, subject, content string) error {
	return n.Bark.Send(key, server, subject, content)
}

func WithDingTalk(token, secret, subject, content string) error {
	return n.DingTalk.Send(token, secret, subject, content)
}

func WithLark(token, secret, subject, content string) error {
	return n.Lark.Send(token, secret, subject, content)
}

func WithFeishu(token, secret, subject, content string) error {
	return n.Lark.SendWithFeishu(token, secret, subject, content)
}

func WithServerChan(userID, sendKey, subject, content string) error {
	return n.ServerChan.Send(userID, sendKey, subject, content)
}
