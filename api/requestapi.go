package api

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/xml"
	"../wx"
)

const (
	token = "farmer"
	EncodingAESKey = "w8rEhj66F7FEFntY76xnxWSw3OJtNGsiPRppBlC8Jsb"
	sCorpID = "wx107a9dfc59f70f80"
)

type wxRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	FromUserName string        `xml:"FromUserName"`
	Content      string        `xml:"Content`
	AgentID      string        `xml:"AgentID"`
}



func ApiRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("Method:", r.Method)
	//It does not work well, update later
	if r.Method == "GET" {
		wechat.ValidateUrl(w, r)
	} else {
		//r.ParseForm()
		//定义
		timestamp := strings.Join(r.Form["timestamp"], "")
		nonce := strings.Join(r.Form["nonce"], "")
		//获取Get传输过来的signature值
		signatureIn := strings.Join(r.Form["msg_signature"], "")
		//打印请求的URL参数
		fmt.Println("path", r.URL.RawQuery)
		fmt.Println(r.Method)
		//fmt.Println(r.Body)
		//validateUrl(w, r)
		body, _ := ioutil.ReadAll(r.Body)
		//sMsg := string(body)
		//fmt.Println(string(body))
		e, err := wechat.NewEncrypter(token, EncodingAESKey, sCorpID)
		if err != nil {
			fmt.Println(err)
			return
		}
		res, _ := e.Decrypt(signatureIn, timestamp, nonce, body)
		//xml_content := string(res)
		//fmt.Println(xml_content)
		//对收到的xml内容做处理
		v := wxRequestBody{}
		err = xml.Unmarshal(res, &v)
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}
		//fmt.Println("FromUserName: ", v.FromUserName)
		fmt.Println("收到的请求内容是: ", v.Content)
		//fmt.Println("AgentID: ", v.AgentID)
		//回复内容的主体，必须是xml格式，有ToUserName，FromUserName,CreateTime,MsgType,Content,MsgID,AgentID
		RequestData := `<xml>
					<ToUserName><![CDATA[` + v.FromUserName + `]]></ToUserName>
					<FromUserName><![CDATA[wx107a9dfc59f70f80]]></FromUserName>
					<CreateTime>` + timestamp + `</CreateTime>
					<MsgType><![CDATA[text]]></MsgType>
					<Content><![CDATA[欢迎您的加入！！！]]></Content>
					<MsgId>7911869847084504763</MsgId>
					<AgentID>` + v.AgentID + `</AgentID>
				</xml>`
		//fmt.Println(RequestData)
		//对回复的xml内容做加密
		ReqValue, _ := e.Encrypt([]byte(RequestData))
		//fmt.Println(string(ReqValue))
		//返回给微信端，微信端会根据加密的内容解密获取content内容
		fmt.Fprintf(w, string(ReqValue))

	}

}