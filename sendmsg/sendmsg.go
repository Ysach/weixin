package sendmsg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
	"time"
)

type SendMsg struct {
	Access_token string `json:"access_token"`
	Current_user_id int  `json:"current_user_id"`
	Key string `json:"key"`
	Phone string `json:"phone"`
	Msg string `json:"msg"`
}

const msg_url = ""

func Sendmsg(msg string, phone string) error {

	var s SendMsg
	s.Access_token = "Ops"
	s.Current_user_id = 123
	s.Key = "^%&&^&GHGH#huhuyhu3yru"
	s.Msg = msg
	s.Phone = phone

	body, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err.Error())
	}

	postReq, err := http.NewRequest("POST", msg_url,bytes.NewReader(body))
	if err != nil {
		fmt.Println(err.Error())
		}
	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{Timeout: 20 * time.Second,}
	presp, err := client.Do(postReq)
	if err != nil {
		return err
	}

	presp.Body.Close()

	return nil
}

