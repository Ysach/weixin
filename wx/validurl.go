package wechat

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"encoding/base64"
	"crypto/aes"
	"crypto/cipher"
	"sort"
	"crypto/sha1"
	"io"
)

const (
	token = "farmer"
	EncodingAESKey = "w8rEhj66F7FEFntY76xnxWSw3OJtNGsiPRppBlC8Jsb"
	sCorpID = "wx107a9dfc59f70f80"
)

func makeSignature(timestamp, nonce string, echostr string) string {
	sl := []string{token, timestamp, nonce, echostr}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))

}

func ValidateUrl(w http.ResponseWriter, r *http.Request) {
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
	res_content_slice := content[:(len(content) - 7)]
	var res_content_slice_new = []uint8{}
	for _, v := range res_content_slice {
		if v != 0 && v != 19 {
			res_content_slice_new = append(res_content_slice_new, v)
		}

	}
	res_content_temp_slice := res_content_slice_new[16:]
	res_content := res_content_temp_slice[:(len(res_content_temp_slice) - len(sCorpID))]
	//获取解密的copID值
	corpID := res_content_temp_slice[(len(res_content_temp_slice) - len(sCorpID)):]
	//fmt.Println(copID)
	fmt.Println(string(corpID))
	if string(corpID) != sCorpID {
		fmt.Println("CorpID不相等，非法的请求!!")
	}

	fmt.Fprintf(w, string(res_content))
}