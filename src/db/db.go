package db

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type UserInfo struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Weight int    `json:"weight"`
}

type Client struct {
	host   string
	port   int
	pwd    string
	db     int
	key    string
	client *redis.Client
}

const keyLock = "lock"
const value = "1"
const expiresLock = 6
const try = 20
const sleepTry = 100

func NewClient(host string, port int, pwd string, db int, key string) *Client {
	c := &Client{}
	c.host = host
	c.port = port
	c.pwd = pwd
	c.db = db
	c.key = key
	return c
}

func (c *Client) Connect() bool {
	addr := c.host + ":" + strconv.Itoa(c.port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: c.pwd,
		DB:       c.db,
	})
	r, err := client.Ping().Result()
	if err != nil {
		log.Println("connect, r", r, ", err:", err)
		return false
	}
	c.client = client
	return true
}

func (c *Client) Disconnect() error {
	return c.client.Close()
}

func (c *Client) AddUser(name string, age int, weight int) bool {
	str, err := json.Marshal(UserInfo{name, age, weight})
	if err != nil {
		log.Println(err)
		return false
	}

	cmd := c.client.HSetNX(c.key, name, str) // 用戶存在則取消新增
	r, err := cmd.Result()
	if err != nil {
		log.Println(err)
		return false
	}
	return r
}

func (c *Client) EditUser(name string, age int, weight int) bool {
	if !c.lock(name) {
		return false
	}

	defer c.unlock(name)

	cmd := c.client.HExists(c.key, name)
	r, err := cmd.Result() // 先檢查 field是否存在, 不存在則離開
	if err != nil {
		log.Println(err)
		return false
	} else if !r {
		log.Println("not exist")
		return false
	}
	str, err := json.Marshal(UserInfo{name, age, weight})
	if err != nil {
		log.Println(err)
		return false
	}

	cmd = c.client.HSet(c.key, name, str)
	_, err = cmd.Result()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (c *Client) DelUser(name string) bool {
	if !c.lock(name) {
		return false
	}

	defer c.unlock(name)

	cmd := c.client.HDel(c.key, name)
	r, err := cmd.Result()
	if err != nil {
		log.Println(err)
		return false
	}
	return r == 1
}

func (c *Client) QueryUser(name string) *UserInfo {
	cmd := c.client.HGet(c.key, name)
	v, err := cmd.Result()
	if err != nil {
		log.Println(err)
		return nil
	}
	obj := &UserInfo{}
	err = json.Unmarshal([]byte(v), obj)
	if err != nil {
		log.Println(err)
		return nil
	}
	return obj
}

func (c *Client) QueryAllUser() []UserInfo {
	cmd := c.client.HVals(c.key)
	r, err := cmd.Result()
	if err != nil {
		log.Println(err)
		return nil
	}
	users := make([]UserInfo, len(r), len(r))
	for i, v := range r {
		obj := &UserInfo{}
		err = json.Unmarshal([]byte(v), obj)
		if err != nil {
			log.Println(err)
			return nil
		}
		users[i] = *obj
	}
	return users
}

func (c *Client) lock(key string) bool {
	lock := false
	for i := 0; i < try; i++ {
		cmd := c.client.SetNX(key+keyLock, value, expiresLock*time.Second)
		r, err := cmd.Result()
		if err != nil {
			log.Println(err)
			time.Sleep(sleepTry * time.Millisecond)
			return false
		}
		if !r {
			time.Sleep(sleepTry * time.Millisecond)
			continue
		}
		lock = true
		break
	}
	return lock
}

func (c *Client) unlock(key string) {
	cmd := c.client.Del(key + keyLock)
	_, err := cmd.Result()
	if err != nil {
		log.Println(err)
	}
}
