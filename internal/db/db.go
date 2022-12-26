package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"hyuga/internal/config"
	"hyuga/internal/db/ent"
	"hyuga/internal/db/ent/predicate"
	"hyuga/internal/db/ent/record"
	"hyuga/internal/db/ent/systemconfig"
	"hyuga/internal/db/ent/user"
	"hyuga/pkg/random"
)

type Client struct {
	*ent.Client
	mux *sync.RWMutex
}

func New(conf *config.Config) (*Client, error) {
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true",
		conf.Db.Username, conf.Db.Password, conf.Db.Address,
		conf.Db.Database, conf.Db.Charset)

	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(100 * time.Second)

	client := &Client{
		Client: ent.NewClient(
			ent.Driver(
				entsql.OpenDB("mysql", db),
			),
		),
		mux: new(sync.RWMutex),
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	if err := client.init(context.Background()); err != nil {
		return nil, err
	}
	return client, err
}

const (
	// 初次部署标识
	InitSystemSetting = "init_system_setting"
)

func (c *Client) init(ctx context.Context) error {
	_, err := c.SystemConfig.Query().Where(systemconfig.KeyEQ(InitSystemSetting)).First(ctx)
	if ent.IsNotFound(err) {
		// 创建默认用户
		b, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = c.User.
			Create().
			SetSid(random.RandomString(5)).
			SetUsername("admin").
			SetPassword(string(b)).
			Save(ctx)
		if err != nil {
			return err
		}
		_, err = c.SystemConfig.Create().
			SetKey(InitSystemSetting).
			SetValue("true").
			Save(ctx)
		if err != nil {
			return err
		}
		return nil

	}
	return err
}

func (c *Client) GetUserByToken(ctx context.Context, token string) (*ent.User, error) {
	return c.User.Query().Where(user.TokenEQ(token)).First(ctx)
}

func (c *Client) GetUserByInviteCode(ctx context.Context, inviteCode string) (*ent.User, error) {
	ic, err := uuid.Parse(inviteCode)
	if err != nil {
		return nil, err
	}
	return c.User.Query().Where(user.InviteCodeEQ(ic)).First(ctx)
}

func (c *Client) GreateUser(ctx context.Context, username, password, inviteCode string) (*ent.User, error) {
	iuser, err := c.GetUserByInviteCode(ctx, inviteCode)
	if err != nil {
		return nil, err
	}
	// 计算hash
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return c.User.
		Create().
		SetSid(random.RandomString(5)).
		SetUsername(username).
		SetPassword(string(b)).
		SetFromUser(iuser.Username).
		Save(ctx)
}

func (c *Client) LoginUser(ctx context.Context, username, password string) (*ent.User, error) {
	u, err := c.User.Query().Where(user.Username(username)).First(ctx)
	if err != nil {
		return nil, err
	}
	// 比较hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, err
	}
	return c.User.UpdateOne(u).
		SetToken(random.RandomString(32)).
		Save(ctx)
}

func (c *Client) LogoutUser(ctx context.Context, uid int) (*ent.User, error) {
	return c.User.UpdateOneID(uid).
		SetToken("").
		Save(ctx)
}

func (c *Client) ChangePasswordUser(ctx context.Context, username, oldpass, newpass string) (*ent.User, error) {
	u, err := c.User.Query().Where(user.Username(username)).First(ctx)
	if err != nil {
		return nil, err
	}
	// 比较hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldpass)); err != nil {
		return nil, err
	}
	// 计算hash
	newpassed, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return c.User.UpdateOne(u).
		SetPassword(string(newpassed)).
		Save(ctx)
}

const (
	DNSType  = "dns"
	HTTPType = "http"
	JndiType = "jndi"
)

func (c *Client) CreateDNSRecord(ctx context.Context, sid, dns, remoteAddr string) (*ent.Record, error) {
	u, err := c.User.Query().Where(user.Sid(sid)).First(ctx)
	if err != nil {
		return nil, err
	}

	return c.Record.Create().
		SetUID(u.ID).
		SetDNSName(dns).
		SetRemoteAddr(remoteAddr).
		SetType(DNSType).
		Save(ctx)
}

func (c *Client) CreateHTTPRecord(ctx context.Context, sid, URL, method, remoteAddr, raw string) (*ent.Record, error) {
	u, err := c.User.Query().Where(user.Sid(sid)).First(ctx)
	if err != nil {
		return nil, err
	}

	return c.Record.Create().
		SetUID(u.ID).
		SetHTTPURL(URL).
		SetHTTPMethod(method).
		SetHTTPRaw(raw).
		SetRemoteAddr(remoteAddr).
		SetType(HTTPType).
		Save(ctx)
}

func (c *Client) CreateJndiRecord(ctx context.Context, sid, protocol, path, remoteAddr string) (*ent.Record, error) {
	u, err := c.User.Query().Where(user.Sid(sid)).First(ctx)
	if err != nil {
		return nil, err
	}

	return c.Record.Create().
		SetUID(u.ID).
		SetJndiProtocol(protocol).
		SetJndiPath(path).
		SetRemoteAddr(remoteAddr).
		SetType(JndiType).
		Save(ctx)
}

func (c *Client) SearchRecord(ctx context.Context, uid int, types, keywords string) ([]*ent.Record, error) {
	wheres := make([]predicate.Record, 0)

	if types != "" {
		wheres = append(wheres, record.TypeEQ(types))
	}
	if keywords != "" {
		switch types {
		case DNSType:
			wheres = append(wheres, record.DNSNameContains(keywords))
		case HTTPType:
			wheres = append(wheres, record.HTTPURLContains(keywords))
			wheres = append(wheres, record.HTTPRawContains(keywords))
		case JndiType:
			wheres = append(wheres, record.JndiPathContainsFold(keywords))
		default:
			wheres = append(wheres, record.DNSNameContains(keywords),
				record.HTTPURLContains(keywords),
				record.HTTPRawContains(keywords),
				record.JndiPathContainsFold(keywords),
			)
		}
	}
	wheres = append(wheres, record.UIDEQ(uid),
		record.CreatedAtGTE(time.Now().Add(-6*time.Hour)))

	return c.Record.
		Query().
		Where(wheres...).
		Order(ent.Desc(record.FieldCreatedAt)).
		Limit(100).All(ctx)
}
