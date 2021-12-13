package database

import (
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	err := Init("redis://localhost:6379/0")
	if err != nil {
		log.Fatal(err)
	}
	m.Run()
}

func TestCreateUser(t *testing.T) {
	user := User{
		ID:    "123",
		Token: "123",
	}
	t.Log(CreateUser(&user))
}

func TestUserExist(t *testing.T) {
	t.Log(UserExist("1123"))
}

func TestSetUserDNSRebinding(t *testing.T) {
	ip := []string{
		"1.1.1.1",
		"2.2.2.2",
	}

	t.Log(SetUserDNSRebinding("123", ip))

	t.Log(GetUserDNSRebinding("123"))
}
