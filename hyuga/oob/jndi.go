package oob

import (
	"bytes"
	"fmt"
	"hyuga/config"
	"hyuga/database"
	"hyuga/oob/limiter"
	"net"
	"strings"
	"time"

	"log"
)

type JndiServer struct {
	addr    string
	closed  chan struct{}
	server  net.Listener
	limiter *limiter.Limiter
}

func NewJndiServer(addr string) *JndiServer {
	server := &JndiServer{
		addr:    addr,
		closed:  make(chan struct{}),
		limiter: limiter.New(10e3),
	}
	return server
}

func (j *JndiServer) ListenAndServe() {
	log.Printf("[jndi] listen on '%s'\n", j.addr)

	listen, err := net.Listen("tcp", j.addr)
	if err != nil {
		log.Printf("[jndi] listen fail error: %s\n", err)
		return
	}
	j.server = listen

LOOP:
	for {
		select {
		case <-j.closed:
			break LOOP
		default:
			j.limiter.Allow()
			conn, err := listen.Accept()
			if err != nil {
				log.Printf("[jndi] listen accept fail error: %s", err)
				j.limiter.Done()
				continue
			}
			go j.acceptProcess(&conn)
		}
	}
}

func (j *JndiServer) Shutdown() error {
	close(j.closed)
	return j.server.Close()
}

func (j *JndiServer) Wait() {
	j.limiter.Wait()
}

/*
thx:
- @4ra1n,@KpLi0rn
- https://4ra1n.love/post/I_AYmmK2J/
*/

func (j *JndiServer) acceptProcess(conn *net.Conn) {
	defer func() {
		(*conn).Close()
		j.limiter.Done()
	}()

	buf := make([]byte, 1024)
	num, err := (*conn).Read(buf)
	if err != nil {
		log.Printf("[jndi] accept data reading err: %s\n", err)
		return
	}
	hexStr := fmt.Sprintf("%x", buf[:num])

	var (
		identity string
		record   database.JndiRecord
	)
	// LDAP Protocol
	if hexStr == ldapfinger {
		if _, err = (*conn).Write(ldapreply); err == nil {
			_, err = (*conn).Read(buf)
			if err != nil {
				log.Printf("[jndi-ldap] read path data err: %s\n", err)
				return
			}
		}

		length := ldapPathLength(buf)
		pathBytes := bytes.Buffer{}
		for i := 1; i <= length; i++ {
			temp := []byte{buf[8+i]}
			pathBytes.Write(temp)
		}

		path := pathBytes.String()
		identity = getSubPath(path)
		if identity != "" {
			record = database.JndiRecord{
				RemoteAddr: (*conn).RemoteAddr().String(),
				Protocol:   "ldap",
				Path:       path,
			}
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

		path := pathBytes.String()
		identity = getSubPath(path)
		if identity != "" {
			record = database.JndiRecord{
				RemoteAddr: (*conn).RemoteAddr().String(),
				Protocol:   "rmi",
				Path:       path,
			}
		}
	}

	if identity != "" && database.UserExist(identity) {
		record.Created = time.Now().Unix()
		if err := database.SetUserRecord(identity, record, config.RecordExpiration); err != nil {
			log.Printf("[jndi] set record '%s' '%#v' error: %s\n", identity, record, err)
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
	if data[0] == 0x4a && data[1] == 0x52 &&
		data[2] == 0x4d && data[3] == 0x49 {
		// version
		if data[4] != 0x00 &&
			data[4] != 0x01 {
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
