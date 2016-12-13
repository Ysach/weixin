package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"crypto/aes"
	"crypto/cipher"
	"io/ioutil"
	"./wx"
	"encoding/xml"
	//"net/url"
)

const (
	token = "farmer"
	EncodingAESKey = "w8rEhj66F7FEFntY76xnxWSw3OJtNGsiPRppBlC8Jsb"
	sCorpID = "wx107a9dfc59f70f80"
)

type wxRequestBody struct {
	XMLName    xml.Name `xml:"xml"`
	FromUserName string	`xml:"FromUserName"`
	Content    string	`xml:"Content`
	AgentID    string	`xml:"AgentID"`
}

func makeSignature(timestamp, nonce string, echostr string) string {
	sl := []string{token, timestamp, nonce, echostr}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))

}

func validateUrl(w http.ResponseWriter, r *http.Request) {
	timestamp := strings.Join(r.Form["timestamp"], "")
	nonce := strings.Join(r.Form["nonce"], "")
	echostr := strings.Join(r.Form["echostr"], "")
	//生产signature值
	signatureGen := makeSignature(timestamp, nonce, echostr)
	//获取Get传输过来的signature值
	signatureIn := strings.Join(r.Form["msg_signature"], "")
	//对比判断，如果相等，则验证通过，明文返回echostr解密的明文msg
	if signatureGen != signatureIn {
		log.Println("不相等")
		return
	}
	//解密收到的加密串儿 echostr

	//先base64解密
	selfKey, err := base64.StdEncoding.DecodeString(EncodingAESKey + "=")
	if err != nil {
		fmt.Println("Error")
	}
	// AES加密
	block, err := aes.NewCipher(selfKey)
	if err != nil {
		fmt.Println("NewChiper Error")
	}
	//CBC mode方式解密
	cfb := cipher.NewCBCDecrypter(block, selfKey[:16])
	content, err := base64.StdEncoding.DecodeString(echostr)
	if err != nil {
		fmt.Println("Content Error")
	}
	//fmt.Println("---content----", content)
	cfb.CryptBlocks(content, content)
	//截取字符串，将明文msg返回给微信
	//res_text := string(content)[20:39]
	res_content_slice := content[:(len(content)-7)]
	var res_content_slice_new = []uint8{}
	for _, v := range res_content_slice {
		if v != 0 && v != 19 {
			res_content_slice_new = append(res_content_slice_new, v)
		}

	}
	res_content_temp_slice := res_content_slice_new[16:]
	res_content := res_content_temp_slice[:(len(res_content_temp_slice)-len(sCorpID))]
	//获取解密的copID值
	corpID := res_content_temp_slice[(len(res_content_temp_slice)-len(sCorpID)):]
	//fmt.Println(copID)
	fmt.Println(string(corpID))
	if string(corpID) != sCorpID {
		fmt.Println("CorpID不相等，非法的请求!!")
	}

	fmt.Fprintf(w, string(res_content))
}

func apiRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("Method:", r.Method)
	//It does not work well, update later
	if r.Method == "GET" {
		validateUrl(w, r)
	}else {
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
		res,_ := e.Decrypt(signatureIn, timestamp, nonce, body)
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
					<ToUserName><![CDATA[`+ v.FromUserName +`]]></ToUserName>
					<FromUserName><![CDATA[wx107a9dfc59f70f80]]></FromUserName>
					<CreateTime>`+ timestamp +`</CreateTime>
					<MsgType><![CDATA[text]]></MsgType>
					<Content><![CDATA[欢迎您的加入！！！]]></Content>
					<MsgId>7911869847084504763</MsgId>
					<AgentID>`+ v.AgentID +`</AgentID>
				</xml>`
		//fmt.Println(RequestData)
		//对回复的xml内容做加密
		ReqValue, _ := e.Encrypt([]byte(RequestData))
		//fmt.Println(string(ReqValue))
		//返回给微信端，微信端会根据加密的内容解密获取content内容
		fmt.Fprintf(w, string(ReqValue))

	}

}

func wxsend(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("Method:", r.Method)
	if r.Method == "POST" {
		msg := strings.Join(r.Form["msg"], "")
		err := wechat.PostCustomMsg("farmer", msg)
		if err != nil {
			log.Println("Post Send Msg error", err)
		}
		log.Println("发送成功！！")
		fmt.Fprintf(w, "发送成功")
	}else {
		log.Println("错误的请求！！")
		fmt.Fprintf(w, "错误的请求！！")
	}
}

//主函数
func main() {
	log.Println("微信服务: 启动!")
	http.HandleFunc("/api/", apiRequest)
	http.HandleFunc("/wx/", wxsend)
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		log.Fatal("微信服务: 端口监听和服务失败！, ", err)
	}
	//服务停止打印log
	log.Println("微信服务: 停止!")

}

