package tool

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	HttpServerHost string `json:"http_server_host"`
	HttpServerPort int    `json:"http_server_port"`
	RedisHost      string `json:"redis_host"`
	RedisPort      int    `json:"redis_port"`
	RedisPwd       string `json:"redis_pwd"`
	RedisDbIndex   int    `json:"redis_dbIndex"`
	RedisKey       string `json:"redis_key"`
}

func ReadConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		log.Println("fail open,", err)
		return nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Println("fail Decode,", err)
		return nil
	}
	return &config
}
