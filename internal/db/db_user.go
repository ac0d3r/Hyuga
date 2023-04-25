package db

import (
	"encoding/json"
	"fmt"

	"github.com/ac0d3r/hyuga/pkg/random"
	"github.com/segmentio/ksuid"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type User struct {
	BaseModel
	GithubUserInfo
	Sid      string `json:"sid"`       // 用于分配子域名
	Token    string `json:"token"`     // token
	APIToken string `json:"api_token"` // API Token
	// TODO
	// ReBing
	// notify
}

type GithubUserInfo struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Home   string `json:"html_url"`
	Avatar string `json:"avatar_url"`
}

func (db *DB) CreateUser(m *User) error {
	if err := db.Create(m); err != nil {
		return err
	}
	return nil
}

func (db *DB) GetUserByGithub(gid int64) (*User, error) {
	m := &User{}
	iter := db.NewIterator(util.BytesPrefix(m.pre()), nil)
	for iter.Next() {
		u := &User{}
		if err := u.decode(iter.Value()); err == nil {
			if u.ID == gid {
				return u, nil
			}
		}
	}

	iter.Release()
	return nil, iter.Error()
}

func (db *DB) GetUserByToken(token string) (*User, error) {
	m := &User{}
	iter := db.NewIterator(util.BytesPrefix(m.pre()), nil)
	for iter.Next() {
		u := &User{}
		if err := u.decode(iter.Value()); err == nil {
			if u.Token == token {
				return u, nil
			}
		}
	}

	iter.Release()
	return nil, iter.Error()
}

func (db *DB) GetUserBySid(sid string) (*User, error) {
	h := &User{Sid: sid}
	if err := db.get(h); err != nil {
		return nil, err
	}
	return h, nil
}

func (db *DB) UserExistWithSid(sid string) bool {
	_, err := db.GetUserBySid(sid)
	return err == nil
}

func (db *DB) UpdateUser(m *User) error {
	if err := db.Update(m); err != nil {
		return err
	}
	return nil
}

var _ Model = (*User)(nil)

func (h *User) id() []byte {
	return append(h.pre(), []byte(h.Sid)...)
}

func (h *User) pre() []byte {
	return []byte(fmt.Sprintf("%s-", modelName(h)))
}

func (h *User) encode() ([]byte, error) {
	return json.Marshal(h)
}

func (h *User) decode(d []byte) error {
	return json.Unmarshal(d, h)
}

func (db *DB) GenUserSid() string {
	var (
		length = 4
		times  = 0
	)

	for {
		sid := random.RandomID(length)
		if !db.UserExistWithSid(sid) {
			return sid
		}
		if times < 3 {
			times++
		} else {
			times = 0
			length++
		}
	}
}

func GenUserToken() string {
	return ksuid.New().String() + random.RandomString(random.RandomInt(5, 20))
}
