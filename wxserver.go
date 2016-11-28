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
)

const (
	token = "farmer"
	EncodingAESKey = "w8rEhj66F7FEFntY76xnxWSw3OJtNGsiPRppBlC8Jsb"
	sCorpID = "wx107a9dfc59f70f80"
	encodingKey = "w8rEhj66F7FEFntY76xnxWSw3OJtNGsiPRppBlC8Jsb="
)

func makeSignature(timestamp, nonce string, echostr string) string {
	sl := []string{token, timestamp, nonce, echostr}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))

}

func validateUrl(w http.ResponseWriter, r *http.Request) bool {
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
		return false
	}
	//解密收到的加密串儿 echostr

	//先base64解密
	selfKey, err := base64.StdEncoding.DecodeString(encodingKey)
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
	res_text := string(content)[20:39]
	//获取解密的copID值
	copID := string(content)[(len(string(content))-len(sCorpID))-7:]
	fmt.Println(copID)
	//fmt.Println(string(content))
	//ln := strings.Trim(string(content), " ")
	//fmt.Println(len(ln))
	//fmt.Println(res_text)
	fmt.Fprintf(w, res_text)
	return true
}

func apiRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !validateUrl(w, r) {
		//log.Println("微信服务：请求非微信平台！")
		//return
	}
	log.Println("微信服务: 微信URL正常！")
}

func main() {
	//aa()
	log.Println("微信服务: 已经启动!")
	http.HandleFunc("/api/", apiRequest)
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		log.Fatal("微信服务: 端口监听和服务失败！, ", err)
	}
	//服务停止打印log
	log.Println("微信服务: 停止!")

}

