package wechat

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	//"strings"
	"strings"
	"bytes"
	//"log"
)

const (
	Corpid = "wx107a9dfc59f70f80"
	Corpsecret = "P1VZn9dDsk4exPAPSbxwGPdnaW2eN-tR6xh3ZDdYAwNRPEgwhkm6k9s1QbVrSzWa"
	Urlsend = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	Urlget = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?"
)

type GetAccessToken struct {
	Access_token string `json:"access_token"`
	Expires_irn int `json:"expires_in"`
}

//定义发送消息的struct
type postMsg struct {
	ToUser string `json:"touser"`
        MsgType string `json:"msgtype"`
	AgentID int `json:"agentid"`
	Safe int `json:"safe"`
        Text PostContent `json:"text"`
}

//定义发送消息的内容struct
type PostContent struct {
	Content string `json:"content"`
}

func PostCustomMsg(toUser, msg string) error  {
	//获取access_token值
	resp, err := http.Get(Urlget + "corpid=" + Corpid + "&" + "corpsecret=" + Corpsecret)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read Failed")
	}
	var m GetAccessToken
	json.Unmarshal(body, &m)
	//fmt.Println(m.Access_token)
	//获取access_token值结束

	pMsg := &postMsg{
		ToUser: toUser,
		MsgType: "text",
		Safe: 0,
		AgentID: 0,
		Text: PostContent{Content: msg},
	}
	pbody, err := json.MarshalIndent(pMsg, " ", " ")
	if err != nil {
		return err
	}
	//fmt.Println(string(pbody))
	postReq, err := http.NewRequest("POST",
		strings.Join([]string{Urlsend, "?access_token=", m.Access_token}, ""),
		bytes.NewReader(pbody))
	if err != nil {
		return err
	}
	//设置Header
	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{}
	presp, err := client.Do(postReq)
	if err != nil {
		return err
	}

	presp.Body.Close()

	return nil

}
/*
func WxSend() {
        //Start Get Token

	resp, err := http.Get(Urlget + "corpid=" + Corpid + "&" + "corpsecret=" + Corpsecret)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read Failed")
	}
	var m GetAccessToken
	json.Unmarshal(body, &m)
	fmt.Println(m.Access_token)
	//End Get Token

	//Start Post Data

	msg := "你好，欢迎！"
	err = PostCustomMsg(m.Access_token, "farmer", msg)
	if err != nil {
		log.Println("Post Send Msg error", err)
	}




}
*/
