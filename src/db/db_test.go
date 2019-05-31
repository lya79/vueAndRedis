package db

import "testing"

const host = "0.0.0.0"
const port = 6379 //6381
const pwd = ""    //"123456789"
const db = 0

func Test_connect(t *testing.T) {
	c := NewClient(host, port, pwd, db, "users")
	success := c.Connect()
	if success == false {
		t.Fatal("err connect")
	}
	err := c.Disconnect()
	if err != nil {
		t.Fatal(err)
	}
}

var client *Client

func Test_init_connect(t *testing.T) {
	client = NewClient(host, port, pwd, db, "users")
	ok := client.Connect()
	if !ok {
		t.Fatal("err connect")
	}

	client.DelUser("aaa")
	client.DelUser("bbb")
}

func Test_add(t *testing.T) {
	ok := client.AddUser("aaa", 11, 111)
	if !ok {
		t.Fatal("err AddUser aaa")
	}

	ok = client.AddUser("bbb", 22, 222)
	if !ok {
		t.Fatal("err AddUser2 bbb")
	}

	ok = client.AddUser("aaa", 22, 222)
	if ok {
		t.Fatal("err AddUser3 ok")
	}
}

func Test_edit(t *testing.T) {
	ok := client.EditUser("aaa", 12, 111)
	if !ok {
		t.Fatal("err new field")
	}
}

func Test_queryUser(t *testing.T) {
	u := client.QueryUser("aaa")
	if u == nil {
		t.Fatal("err u == nil")
	}
	t.Log(u)
}

func Test_queryAllUser(t *testing.T) {
	u := client.QueryAllUser()
	if u == nil {
		t.Fatal("err u == nil")
	}
	for _, v := range u {
		t.Log(v)
	}
}
func Test_disconnect(t *testing.T) {
	err := client.Disconnect()
	if err != nil {
		t.Fatal(err)
	}
}
