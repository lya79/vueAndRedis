package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"vue/src/db"
	"vue/src/tool"
)

type Handler struct {
	Config tool.Config

	client *db.Client
}

type RespJSON struct {
	ErrCode int         `json:"error-code"`
	ErrText string      `json:"error-text"`
	Data    interface{} `json:"data"`
}

var createRegexp = regexp.MustCompile(`^/create/user/{0,1}$`)
var updateRegexp = regexp.MustCompile(`^/update/user/{0,1}$`)
var delRegexp = regexp.MustCompile(`^/del/user/{0,1}$`)
var queryRegexp = regexp.MustCompile(`^/query/user/{0,1}$`)

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req)

	matchURL := h.createUser(w, req) || h.updateUser(w, req) || h.delUser(w, req) || h.queryUser(w, req)
	if matchURL {
		return
	}

	if req.URL.Path == "/" || req.URL.Path == "/index.html" || req.URL.Path == "/index.css" || req.URL.Path == "/index.js" {
		http.ServeFile(w, req, "./static/"+req.URL.Path[1:])
		return
	}

	http.ServeFile(w, req, "./static/404.jpg")
}

func (h *Handler) Connect() bool {
	host := h.Config.RedisHost
	port := h.Config.RedisPort
	pwd := h.Config.RedisPwd
	dbIndex := h.Config.RedisDbIndex
	key := h.Config.RedisKey

	h.client = db.NewClient(host, port, pwd, dbIndex, key)

	ok := h.client.Connect()
	if !ok {
		log.Println("fail connect redis")
	} else {
		log.Println("success connect redis")
	}
	return ok
}

func (h *Handler) Disonnect() {
	err := h.client.Disconnect()
	if err != nil {
		log.Println("fail disonnect redis, ", err)
		return
	}

	log.Println("success disonnect redis")
}

func (h *Handler) createUser(w http.ResponseWriter, req *http.Request) bool {
	if !createRegexp.Match([]byte(req.URL.Path)) {
		log.Println("not match createRegexp,", req.URL.Path)
		return false
	}

	if req.Method != "POST" {
		http.NotFound(w, req)
		return true
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var u db.UserInfo

	contentType := req.Header.Get("Content-type")
	if contentType == "application/json" {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&u)
		if err != nil {
			log.Println("fail Decode,", req.Body)
			respData, err := json.Marshal(RespJSON{ErrCode: 3, ErrText: "參數錯誤", Data: ""})
			if err != nil {
				log.Println("fail Marshal, req:", req, ", ", err)
			} else {
				w.Write(respData)
			}
			return true
		}
		log.Println("success Decode,", u)
	} else {
		req.ParseForm()
		name := req.Form["name"][0]
		age, err := strconv.Atoi(req.Form["age"][0])
		weight, err2 := strconv.Atoi(req.Form["weight"][0])
		if err != nil || err2 != nil {
			respData, err := json.Marshal(RespJSON{ErrCode: 3, ErrText: "參數錯誤", Data: ""})
			if err != nil {
				log.Println("fail Marshal, req:", req, ", ", err)
			} else {
				w.Write(respData)
			}
			return true
		}
		u.Name = name
		u.Age = age
		u.Weight = weight
	}

	ok := h.client.AddUser(u.Name, u.Age, u.Weight)
	if !ok {
		respData, err := json.Marshal(RespJSON{ErrCode: 1, ErrText: "用戶已經存在", Data: ""})
		if err != nil {
			log.Println("fail Marshal, req:", req, ", ", err)
		} else {
			w.Write(respData)
		}
		return true
	}

	respData, err := json.Marshal(RespJSON{ErrCode: 0, ErrText: "", Data: ""})
	if err != nil {
		log.Println("fail Marshal, req:", req, ", ", err)
	} else {
		w.Write(respData)
	}

	return true
}

func (h *Handler) updateUser(w http.ResponseWriter, req *http.Request) bool {
	if !updateRegexp.Match([]byte(req.URL.Path)) {
		log.Println("not match updateRegexp,", req.URL.Path)
		return false
	}

	if req.Method != "POST" {
		http.NotFound(w, req)
		return true
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	q := req.URL.Query()
	name := q.Get("name")

	var u db.UserInfo
	u.Name = name

	contentType := req.Header.Get("Content-type")
	if contentType == "application/json" {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&u)
		if err != nil {
			log.Println("fail Decode,", req.Body)
			respData, err := json.Marshal(RespJSON{ErrCode: 3, ErrText: "參數錯誤", Data: ""})
			if err != nil {
				log.Println("fail Marshal, req:", req, ", ", err)
			} else {
				w.Write(respData)
			}
			return true
		}
		log.Println("success Decode,", u)
	} else {
		req.ParseForm()
		age2, err := strconv.Atoi(req.Form["age"][0])
		weight2, err2 := strconv.Atoi(req.Form["weight"][0])
		if err != nil || err2 != nil {
			respData, err := json.Marshal(RespJSON{ErrCode: 3, ErrText: "參數錯誤", Data: ""})
			if err != nil {
				log.Println("fail Marshal, req:", req, ", ", err)
			} else {
				w.Write(respData)
			}
			return true
		}
		u.Age = age2
		u.Weight = weight2
	}

	ok := h.client.EditUser(u.Name, u.Age, u.Weight)
	if !ok {
		respData, err := json.Marshal(RespJSON{ErrCode: 2, ErrText: "用戶不存在", Data: ""})
		if err != nil {
			log.Println("fail Marshal, req:", req, ", ", err)
		} else {
			w.Write(respData)
		}
		return true
	}

	respData, err := json.Marshal(RespJSON{ErrCode: 0, ErrText: "", Data: ""})
	if err != nil {
		log.Println("fail Marshal, req:", req, ", ", err)
	} else {
		w.Write(respData)
	}

	return true
}

func (h *Handler) delUser(w http.ResponseWriter, req *http.Request) bool {
	if !delRegexp.Match([]byte(req.URL.Path)) {
		log.Println("not match delRegexp,", req.URL.Path)
		return false
	}

	if req.Method != "POST" {
		http.NotFound(w, req)
		return true
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	q := req.URL.Query()
	name := q.Get("name")

	ok := h.client.DelUser(name)
	if !ok {
		respData, err := json.Marshal(RespJSON{ErrCode: 2, ErrText: "用戶不存在", Data: ""})
		if err != nil {
			log.Println("fail Marshal, req:", req, ", ", err)
		} else {
			w.Write(respData)
		}
		return true
	}

	respData, err := json.Marshal(RespJSON{ErrCode: 0, ErrText: "", Data: ""})
	if err != nil {
		log.Println("fail Marshal, req:", req, ", ", err)
	} else {
		w.Write(respData)
	}

	return true
}

func (h *Handler) queryUser(w http.ResponseWriter, req *http.Request) bool {
	if !queryRegexp.Match([]byte(req.URL.Path)) {
		log.Println("not match queryRegexp,", req.URL.Path)
		return false
	}

	if req.Method != "GET" {
		http.NotFound(w, req)
		return true
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var data interface{}

	q := req.URL.Query()
	name := q.Get("name")

	if name != "" {
		userinfo := h.client.QueryUser(name)
		if userinfo == nil {
			data3, err := json.Marshal(RespJSON{ErrCode: 2, ErrText: name + "用戶不存在", Data: ""})
			if err != nil {
				log.Println("fail Marshal, req:", req, ", ", err)
			} else {
				w.Write(data3)
			}
			return true
		}
		data = userinfo
	} else {
		userinfoSlice := h.client.QueryAllUser()
		if userinfoSlice == nil {
			data3, err := json.Marshal(RespJSON{ErrCode: 3, ErrText: "系統錯誤", Data: ""})
			if err != nil {
				log.Println("fail Marshal, req:", req, ", ", err)
			} else {
				w.Write(data3)
			}
			return true
		}
		data = userinfoSlice
	}

	respData, err := json.Marshal(RespJSON{ErrCode: 0, ErrText: "", Data: data})
	if err != nil {
		log.Println("fail Marshal, req:", req, ", ", err)
	} else {
		w.Write(respData)
	}

	log.Println("success queryUser")

	return true
}
