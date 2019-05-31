package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"vue/src/handler"
	"vue/src/tool"
)

const pathOfConfig = "./config.json"

func main() {
	log.Println("server running")
	defer log.Println("server stopped")

	config := tool.ReadConfig(pathOfConfig)
	if config == nil {
		log.Println("fail ReadConfig, path:", pathOfConfig)
		return
	}

	c, err := json.Marshal(config)
	if err != nil {
		log.Println("fail json.Marshal(userinfo), ", err)
		return
	}
	log.Println(string(c))

	addr := config.HttpServerHost + ":" + strconv.Itoa(config.HttpServerPort)

	mux, handler := initServeMux(*config)
	if mux == nil {
		log.Println("fail initServeMux")
		return
	} else {
		log.Println("success initServeMux")
	}
	defer handler.Disonnect()

	serv := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	err = serv.ListenAndServe()
	if err != nil {
		log.Println("ListenAndServe, error, ", err)
	}
}

func initServeMux(config tool.Config) (*http.ServeMux, *handler.Handler) {
	mux := http.NewServeMux()

	h := &handler.Handler{Config: config}
	ok := h.Connect()
	if !ok {
		log.Println("fail connect redis")
		return nil, nil
	}

	mux.Handle("/", h)

	return mux, h
}
