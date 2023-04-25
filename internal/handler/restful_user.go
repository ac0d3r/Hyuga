package handler

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ac0d3r/hyuga/internal/db"
	"github.com/gin-gonic/gin"
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

	// github user info
	accessToken := payload["access_token"]
	info, err := getUserInfo(accessToken)
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

	// TODO
	c.SetCookie("sid", u.Sid, 0, "/", "", false, true)
	c.SetCookie("token", u.Token, 0, "/", "", false, true)
	c.Redirect(http.StatusFound, "/api/v2/user/info")
}

func getUserInfo(token string) (*db.GithubUserInfo, error) {
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

}

func (w *restfulHandler) reset_token(c *gin.Context) {

}

func (w *restfulHandler) logout(c *gin.Context) {

}
