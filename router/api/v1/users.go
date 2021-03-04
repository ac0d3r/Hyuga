package v1

import (
	"context"
	"fmt"
	"net/http"

	"hyuga/conf"
	"hyuga/db"
	"hyuga/internal"
	"hyuga/internal/random"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	ctx context.Context = context.Background()
)

func genUser() (identity string, token string, err error) {
	// remove UpperCharset
	c := random.Charset + random.Number
	length := 4
	times := 0

	for {
		identity = random.StringWithCharset(length, c)
		exist, err := db.IdentityExist(ctx, identity)
		if err != nil {
			return "", "", err
		}
		if !exist {
			token = random.String(random.Int(32, 64))
			exist, err := db.UserExist(ctx, identity, token)
			if err != nil {
				return "", "", err
			}
			if !exist {
				break
			}
		}
		if times < 3 {
			times++
		} else {
			times = 0
			length++
		}
	}
	return
}

// CreateUser create user
func CreateUser(c echo.Context) error {
	identity, token, err := genUser()
	if err != nil {
		return c.JSON(internal.ProcessingError(err))
	}
	log.Debug("api/v1/CreateUser ", identity, " / ", token)
	if err = db.CreateUser(ctx, identity, token); err != nil {
		return c.JSON(internal.ProcessingError(err))
	}

	return c.JSON(http.StatusOK, &internal.RespJSON{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    map[string]string{"identity": fmt.Sprintf("%s.%s", identity, conf.Domain.Main), "token": token},
		Success: true,
	})
}

type setRebinding struct {
	Token string   `json:"token"`
	Hosts []string `json:"hosts"`
}

// GetUserDNSRebinding get user dns-rebinding hosts
func GetUserDNSRebinding(c echo.Context) error {
	identity := c.Param("identity")
	token := c.QueryParam("token")

	log.Debug("api/v1/GetUserDNSRebinding ", identity, token)
	exist, err := db.UserExist(ctx, identity, token)
	if err != nil {
		return c.JSON(internal.ProcessingError(err))
	}
	if !exist {
		code := http.StatusUnauthorized
		return c.JSON(code, &internal.RespJSON{
			Code:    code,
			Message: http.StatusText(code),
			Data:    nil,
			Success: false,
		})
	}
	ips, err := db.GetUserDNSRebinding(ctx, identity, false)
	if err != nil {
		return c.JSON(internal.ProcessingError(err))
	}

	code := http.StatusOK
	return c.JSON(code, &internal.RespJSON{
		Code:    code,
		Message: http.StatusText(code),
		Data:    map[string][]string{"rebinding_hosts": ips},
		Success: true,
	})
}

// SetUserDNSRebinding set user dns-rebinding
func SetUserDNSRebinding(c echo.Context) error {
	identity := c.Param("identity")
	dnsRebinding := setRebinding{}

	if err := c.Bind(&dnsRebinding); err != nil {
		code, resp := internal.ProcessingError(err)
		return c.JSON(code, resp)
	}
	// check dnsrebinging ip
	for index, ip := range dnsRebinding.Hosts {
		ok, err := internal.CheckIP(ip)
		if err != nil || !ok {
			return c.JSON(internal.ProcessingError(fmt.Errorf(`Invalid Parameter 'hosts[%d]' "%s"`, index, ip)))
		}
	}
	log.Debug("api/v1/SetUserDNSRebinding ", identity, dnsRebinding.Hosts)

	err := db.SetUserDNSRebinding(ctx, identity,
		dnsRebinding.Token,
		internal.StringSlice2AnySlice(dnsRebinding.Hosts))
	if err != nil {
		code, resp := internal.ProcessingError(err)
		return c.JSON(code, resp)
	}

	code := http.StatusOK
	return c.JSON(code, &internal.RespJSON{
		Code:    code,
		Message: http.StatusText(code),
		Data:    nil,
		Success: true,
	})
}
