package tool

import (
	"encoding/json"
	"testing"
)

func Test_rw(t *testing.T) {
	config := ReadConfig("../config.json")
	t.Log(config)

	if config == nil {
		t.Fatal("fail readconfig")
	}

	c, err := json.Marshal(config)
	if err != nil {
		t.Fatal("fail json.Marshal(userinfo), ", err)
		return
	}

	t.Log(string(c))
}
