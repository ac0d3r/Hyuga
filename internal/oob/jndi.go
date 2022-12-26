package oob

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strings"

	"hyuga/internal/config"
	"hyuga/internal/db"
	"hyuga/pkg/limiter"

	"github.com/sirupsen/logrus"
)

type Jndi struct {
	cnf *config.OOB
	db  *db.Client

	closed chan struct{}
	s      net.Listener
	l      *limiter.Limiter
}

func NewJndi(cnf *config.OOB, db *db.Client) *Jndi {
	return &Jndi{
		cnf: cnf,
		db:  db,

		closed: make(chan struct{}),
		l:      limiter.New(10e3),
	}
}

func (j *Jndi) ListenAndServe() error {
	var err error
	logrus.Infof("[jndi] listen on '%s'", j.cnf.Jndi.Address)
	j.s, err = net.Listen("tcp", j.cnf.Jndi.Address)
	if err != nil {
		logrus.Warnf("[jndi] listen fail error: %s", err)
		return err
	}

LOOP:
	for {
		select {
		case <-j.closed:
			break LOOP
		default:
			j.l.Allow()
			conn, err := j.s.Accept()
			if err != nil {
				logrus.Warnf("[jndi] listen accept fail error: %s", err)
				j.l.Done()
				continue
			}
			go j.acceptProcess(&conn)
		}
	}
	return nil
}

func (j *Jndi) Shutdown() error {
	close(j.closed)
	defer j.l.Wait()

	return j.s.Close()
}

/*
thx:
- @4ra1n,@KpLi0rn
- https://4ra1n.love/post/I_AYmmK2J/
*/

func (j *Jndi) acceptProcess(conn *net.Conn) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func() {
		(*conn).Close()
		j.l.Done()
	}()

	buf := make([]byte, 1024)
	num, err := (*conn).Read(buf)
	if err != nil {
		logrus.Warnf("[jndi] accept data reading err: %s", err)
		return
	}
	hexStr := fmt.Sprintf("%x", buf[:num])

	var (
		sid      string
		protocol string
		path     string
	)
	// LDAP Protocol
	if hexStr == ldapfinger {
		if _, err = (*conn).Write(ldapreply); err == nil {
			_, err = (*conn).Read(buf)
			if err != nil {
				logrus.Warnf("[jndi-ldap] read path data err: %s", err)
				return
			}
		}

		length := ldapPathLength(buf)
		pathBytes := bytes.Buffer{}
		for i := 1; i <= length; i++ {
			temp := []byte{buf[8+i]}
			pathBytes.Write(temp)
		}

		path = pathBytes.String()
		sid = getSubPath(path)
		if sid != "" {
			protocol = "ldap"
		}
	}

	// RMI Protocol
	if checkRMI(buf) {
		_, _ = (*conn).Write(rmireplay)
		// 这里读到的数据没有用处
		_, _ = (*conn).Read(buf)
		// 需要发一次空数据然后接收call信息
		_, _ = (*conn).Write([]byte{})
		_, _ = (*conn).Read(buf)

		var dataList []byte
		var flag bool
		// 从后往前读因为空都是00
		for i := len(buf) - 1; i >= 0; i-- {
			// 这里要用一个flag来区分
			// 因为正常数据中也会含有00
			if buf[i] != 0x00 || flag {
				flag = true
				dataList = append(dataList, buf[i])
			}
		}
		// 已读到的长度等于当前读到的字节代表的数字
		// 那么认为已读到的字符串翻转后是路径参数
		var j int
		for i := 0; i < len(dataList); i++ {
			if int(dataList[i]) == i {
				j = i
				break
			}
		}

		if len(dataList) < j {
			return
		}
		temp := dataList[0:j]
		pathBytes := &bytes.Buffer{}
		// 翻转后拿到真正的路径参数
		for i := len(temp) - 1; i >= 0; i-- {
			pathBytes.Write([]byte{dataList[i]})
		}

		path = pathBytes.String()
		sid = getSubPath(path)
		if sid != "" {
			protocol = "rmi"
		}
	}

	if sid != "" {
		remoteAddr := strings.Split((*conn).RemoteAddr().String(), ":")[0]
		if _, err := j.db.CreateJndiRecord(ctx, sid, protocol, path, remoteAddr); err != nil {
			logrus.Warnf("[jndi] set record '%s' %s/%s error: %s", sid, protocol, path, err)
		}
	}
}

func getSubPath(s string) string {
	i := strings.Index(strings.TrimLeft(s, "/"), "/")
	if i <= 0 {
		return ""
	}
	return s[:i]
}

var (
	// ldap protocol
	// https://ldap.com/ldapv3-wire-protocol-reference-bind/
	/*
		30 0c -- Begin the LDAPMessage sequence
			02 01 01 --  The message ID (integer value 1)
			60 07 -- Begin the bind request protocol op
				02 01 03 -- The LDAP protocol version (integer value 3)
				04 00 -- Empty bind DN (0-byte octet string)
				80 00 -- Empty password (0-byte octet string with type context-specific
					-- primitive zero)
	*/
	ldapfinger string = "300c020101600702010304008000"
	/*
		30 0c -- Begin the LDAPMessage sequence
			02 01 01 -- The message ID (integer value 1)
			61 07 -- Begin the bind response protocol op
				0a 01 00 -- success result code (enumerated value 0)
				04 00 -- No matched DN (0-byte octet string)
				04 00 -- No diagnostic message (0-byte octet string)
	*/
	ldapreply = []byte{
		0x30, 0x0c,
		0x02, 0x01, 0x01,
		0x61, 0x07,
		0x0a, 0x01, 0x00,
		0x04, 0x00,
		0x04, 0x00,
	}
)

func ldapPathLength(buf []byte) int {
	if len(buf) < 9 {
		return 0
	}
	length := buf[8]
	if len(buf) < 9+int(length) {
		return 0
	}
	return int(length)
}

var (
	// rmi protocol
	// https://docs.oracle.com/javase/9/docs/specs/rmi/protocol.html
	rmireplay = []byte{
		0x4e, 0x00, 0x09, // 保证4e00开头
		0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, // 模拟 127.0.0.1
		0x00, 0x00, 0xc4, 0x12,
	}
)

func checkRMI(data []byte) bool {
	if len(data) < 8 {
		return false
	}
	// header
	if data[0] == 0x4a &&
		data[1] == 0x52 &&
		data[2] == 0x4d &&
		data[3] == 0x49 {
		// version
		if data[4] != 0x00 {
			return false
		}
		if data[5] != 0x01 &&
			data[5] != 0x02 {
			return false
		}

		// protocol
		if data[6] != 0x4b &&
			data[6] != 0x4c &&
			data[6] != 0x4d {
			return false
		}
		lastData := data[7:]
		for _, v := range lastData {
			if v != 0x00 {
				return false
			}
		}
		return true
	}

	return false
}
