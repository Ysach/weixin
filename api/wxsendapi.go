package api

import (
	"net/http"
	"strings"
	"runtime"
	"fmt"
	"log"
	"encoding/json"
	"../db/mysql"
	"../sendmsg"
	"../wx"
	"time"
)

type HttpRender struct {
	Msg string `json:"msg"`
	Status int `json:"status"`
}

func Wxsend(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// fmt.Println("Method:", r.Method)
	getType := r.Header.Get("Content-Type")
	var detail HttpRender
	//fmt.Println(r.Header.Get("Content-Type"))
	//fmt.Println(r.Header)
	switch r.Method {
	case "POST":
		if getType != "application/x-www-form-urlencoded" {
			http.Error(w, "请使用application/x-www-form-urlencoded", 404)
			return
		}
		msg := strings.Join(r.Form["msg"], "")
		runtime.GOMAXPROCS(runtime.NumCPU())
		go wechat.PostCustomMsg("farmer", msg)
		//erro := sendmsg.Sendmsg(msg, phone)
		//if erro != nil {
		//	fmt.Println(erro.Error())
		//}
		err := mysqldb.MysqlData(msg)
		s := mysqldb.MysqlSel()
		fmt.Println("Start Send msg", time.Now().Unix())
		erro := sendmsg.Sendmsg(msg, s)
		if erro != nil {
			fmt.Println(erro.Error())
		}
		fmt.Println("End Send msg", time.Now().Unix())


		//fmt.Println(p)
		if err != nil {
			fmt.Fprintf(w, "数据库存储失败！！")
		}

		log.Println("发送成功！！")
		detail.Msg = "success"
		detail.Status = http.StatusOK
		body, err := json.MarshalIndent(detail, " ", " ")
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Fprintf(w, string(body))

	case "GET":
		http.Error(w, "BAD Request Method", 405)
		log.Println("错误的请求！！---", r.Method, r.RequestURI)
	}

}
