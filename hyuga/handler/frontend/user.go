package frontend

import (
	"fmt"
	"hyuga/config"
	"hyuga/database"
	"hyuga/handler/base"
	"hyuga/handler/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

func CreateUser(c *gin.Context) {
	user := &database.User{
		ID:      genID(),
		Token:   ksuid.New().String() + util.RandomString(util.RandomInt(5, 20)),
		Created: time.Now().Unix(),
	}
	if err := database.CreateUser(user); err != nil {
		base.ReturnError(c, 102, err)
		return
	}
	base.ReturnJSON(c, map[string]string{
		"id":    fmt.Sprintf("%s.%s", user.ID, config.MainDomain),
		"token": user.Token,
	})
}

func DeleteUser(c *gin.Context) {
	if err := database.DeleteUserByUserID(base.GetUserID(c)); err != nil {
		base.ReturnError(c, 102, err)
		return
	}
	records := []database.Record{
		database.DnsRecord{},
		database.HttpRecord{},
	}
	for _, r := range records {
		if err := database.DeleteRecordsByUserID(r, base.GetUserID(c)); err != nil {
			base.ReturnError(c, 102, err)
			return
		}
	}
	base.ReturnJSON(c, nil)
}

func GetUserDnsRebinding(c *gin.Context) {
	ips, err := database.GetUserDNSRebinding(base.GetUserID(c))
	if err != nil {
		base.ReturnError(c, 102, err)
		return
	}
	base.ReturnJSON(c, ips)
}

func UpdateUserDnsRebinding(c *gin.Context) {
	param := &struct {
		IP []string `json:"ip,omitempty" form:"ip" binding:"required"`
	}{}
	if err := c.ShouldBind(param); err != nil {
		base.ReturnError(c, 101, err)
		return
	}
	if err := database.SetUserDNSRebinding(base.GetUserID(c), param.IP); err != nil {
		base.ReturnError(c, 102, err)
		return
	}

	base.ReturnJSON(c, nil)
}

// genID to get a short ID
func genID() string {
	var (
		length = 4
		times  = 0
	)

	for {
		id := util.RandomID(length)
		if !database.UserExist(id) {
			return id
		}
		if times < 3 {
			times++
		} else {
			times = 0
			length++
		}
	}
}
