package handler

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/ac0d3r/hyuga/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	githubOauth = "https://github.com/login/oauth/access_token"
	githubUser  = "https://api.github.com/user"
)

var (
	httpc *http.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

func (w *restfulHandler) login(c *gin.Context) {
	type Request struct {
		Code string `json:"code" form:"code" binding:"required"`
	}

	param := new(Request)
	if BindParam(c, param) {
		return
	}
	logrus.Infof("[restful] user login.")

	data, err := json.Marshal(map[string]string{
		"client_id":     w.cnf.Github.ClientID,
		"client_secret": w.cnf.Github.ClientSecret,
		"code":          param.Code,
	})
	if err != nil {
		ReturnError(c, errOauth, err)
		return
	}

	req, err := http.NewRequest("POST", githubOauth, bytes.NewReader(data))
	if err != nil {
		ReturnError(c, errOauth, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := httpc.Do(req)
	if err != nil {
		ReturnError(c, errOauth, err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		ReturnError(c, errOauth, err)
		return
	}
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		ReturnError(c, errOauth, err)
		return
	}

	payload := make(map[string]string)
	if err := json.Unmarshal(data, &payload); err != nil {
		ReturnError(c, errOauth, err)
		return
	}

	logrus.Infof("[restful] authorization successful and get user info.")
	// github user info
	accessToken := payload["access_token"]
	info, err := getGithubUserInfo(accessToken)
	if err != nil {
		ReturnError(c, errOauth, err)
		return
	}
	// get github user from db
	u, err := w.db.GetUserByGithub(info.ID)
	if err != nil {
		ReturnError(c, errDatabase, err)
		return
	}
	if u == nil {
		// not found, create new user
		u = &db.User{
			GithubUserInfo: *info,
			Sid:            w.db.GenUserSid(),
			Token:          db.GenUserToken(),
			APIToken:       db.GenUserToken(),
		}
		if err := w.db.CreateUser(u); err != nil {
			ReturnError(c, errDatabase, err)
			return
		}
	} else {
		// found, update access token
		u.GithubUserInfo = *info
		u.Token = db.GenUserToken()
		if err := w.db.UpdateUser(u); err != nil {
			ReturnError(c, errDatabase, err)
			return
		}
	}

	logrus.Infof("[restful] login success.")
	setCookie(c, "sid", u.Sid)
	setCookie(c, "token", u.Token)
	c.Redirect(http.StatusFound, "/")
}

func getGithubUserInfo(token string) (*db.GithubUserInfo, error) {
	req, err := http.NewRequest("GET", githubUser, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("get github user info failed")
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	info := &db.GithubUserInfo{}
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	return info, nil
}

func (w *restfulHandler) info(c *gin.Context) {
	user, err := w.db.GetUserBySid(c.GetString("sid"))
	if err != nil {
		ReturnError(c, errDatabase, err)
		return
	}

	_, port, _ := net.SplitHostPort(w.oob.JNDI.Address)
	ReturnJSON(c, map[string]any{
		"name":      user.Name,
		"avatar":    user.Avatar,
		"sid":       user.Sid,
		"token":     user.APIToken,
		"rebinding": user.DnsRebind.DNS,
		"data": map[string]string{
			"subdomain": fmt.Sprintf("%s.%s", user.Sid, w.oob.DNS.Main),
			"rdomain":   fmt.Sprintf("r.%s.%s", user.Sid, w.oob.DNS.Main),
			"ldap":      fmt.Sprintf("ldap://%s:%s/%s/", w.oob.DNS.Main, port, user.Sid),
			"rmi":       fmt.Sprintf("rmi://%s:%s/%s/", w.oob.DNS.Main, port, user.Sid),
		},
		"notify": user.Notify,
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

func (w *restfulHandler) record(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Warnf("[restful] upgrade websocket failed: %v", err)
		return
	}

	sid := c.GetString("sid")
	go func() {
		defer ws.Close()

		logrus.Infof("[restful] start record stream")
		c := make(chan struct{})
		go func() {
			if _, _, err := ws.ReadMessage(); err != nil {
				close(c)
			}
		}()
		// get user all records
		records, err := w.recorder.Get(sid)
		if err != nil {
			logrus.Warnf("[restful] get user records err: %s", err.Error())
			return
		}
		for _, r := range records {
			if err = ws.WriteJSON(r); err != nil {
				logrus.Infof("[restful][stream] push record err: %s", err.Error())
			}
		}
		// subscribe user record event
		s := w.eventbus.Subscribe(sid)
		defer w.eventbus.Unsubscribe(s)
		for {
			select {
			case <-c:
				logrus.Infof("[restful] close record stream")
				return
			case msg := <-s.Out():
				logrus.Infof("[restful][stream] push record msg: %v", msg)
				if err = ws.WriteJSON(msg); err != nil {
					return
				}
			}
		}
	}()
}

func (w *restfulHandler) notify(c *gin.Context) {
	type Request struct {
		Enable bool `json:"enable" form:"enable"`
		Bark   struct {
			Key    string `json:"key" form:"key"`
			Server string `json:"server" form:"server"`
		} `json:"bark" form:"bark"`
		Dingtalk struct {
			Token  string `json:"token" form:"token"`
			Secret string `json:"secret" form:"secret"`
		} `json:"dingtalk" form:"dingtalk"`
		Lark struct {
			Token  string `json:"token" form:"token"`
			Secret string `json:"secret" form:"secret"`
		} `json:"lark" form:"lark"`
		Feishu struct {
			Token  string `json:"token" form:"token"`
			Secret string `json:"secret" form:"secret"`
		} `json:"feishu" form:"feishu"`
		ServerChan struct {
			UserID  string `json:"user_id" form:"user_id"`
			SendKey string `json:"send_key" form:"send_key"`
		} `json:"server_chan" form:"server_chan"`
	}

	param := new(Request)
	if BindParam(c, param) {
		return
	}

	user, err := w.db.GetUserBySid(c.GetString("sid"))
	if err != nil {
		ReturnError(c, errDatabase, err)
		return
	}
	user.Notify.Enable = param.Enable
	user.Notify.Bark.Key = param.Bark.Key
	user.Notify.Bark.Server = param.Bark.Server
	user.Notify.Dingtalk.Token = param.Dingtalk.Token
	user.Notify.Dingtalk.Secret = param.Dingtalk.Secret
	user.Notify.Lark.Token = param.Lark.Token
	user.Notify.Lark.Secret = param.Lark.Secret
	user.Notify.Feishu.Token = param.Feishu.Token
	user.Notify.Feishu.Secret = param.Feishu.Secret
	user.Notify.ServerChan.UserID = param.ServerChan.UserID
	user.Notify.ServerChan.SendKey = param.ServerChan.SendKey

	if err := w.db.UpdateUser(user); err != nil {
		ReturnError(c, errDatabase, err)
		return
	}
	ReturnJSON(c, nil)
}

func (w *restfulHandler) rebinding(c *gin.Context) {
	type Request struct {
		Rebinding []string `json:"rebinding" form:"rebinding" binding:"required"`
	}
	param := new(Request)
	if BindParam(c, param) {
		return
	}

	user, err := w.db.GetUserBySid(c.GetString("sid"))
	if err != nil {
		ReturnError(c, errDatabase, err)
		return
	}
	user.DnsRebind.DNS = param.Rebinding
	if err := w.db.UpdateUser(user); err != nil {
		ReturnError(c, errDatabase, err)
		return
	}
	ReturnJSON(c, nil)
}
func (w *restfulHandler) reset(c *gin.Context) {
	user, err := w.db.GetUserBySid(c.GetString("sid"))
	if err != nil {
		ReturnError(c, errDatabase, err)
		return
	}

	user.APIToken = db.GenUserToken()
	if err := w.db.UpdateUser(user); err != nil {
		ReturnError(c, errDatabase, err)
		return
	}

	ReturnJSON(c, nil)
}

func (w *restfulHandler) logout(c *gin.Context) {
	user, err := w.db.GetUserBySid(c.GetString("sid"))
	if err != nil {
		ReturnError(c, errDatabase, err)
		return
	}

	user.Token = ""
	if err := w.db.UpdateUser(user); err != nil {
		ReturnError(c, errDatabase, err)
		return
	}

	ReturnJSON(c, nil)
}
