package api

import (
	"net/http"
	//"strings"
	//"runtime"
	"fmt"
	"log"
	"encoding/json"
	"../db/mysql"
	"../email"
	"../sendmsg"
	"../wx"
	"time"
	"net/url"
)

type HttpRender struct {
	Msg string `json:"msg"`
	Status int `json:"status"`
}

type DataPasser struct {

	FormData chan url.Values
}

func Wxsend(w http.ResponseWriter, r *http.Request) { //不使用channel
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
		//msg := strings.Join(r.Form["msg"], "")
		msg := r.Form.Get("msg")
		// runtime.GOMAXPROCS(runtime.NumCPU())

		err := wechat.PostCustomMsg("farmer", msg) // Send wechat messages
		if err != nil {
			fmt.Println("Send wechat msg error")
		}
		MysqlErr := mysqldb.MysqlData(msg)	// save msg in db
		if MysqlErr != nil {
			fmt.Println("Save db Error")
		}
		s, e := mysqldb.MysqlSel()	// Get phone and email string from db
		var email sendemail.Email

		email.To = e
		email.Mailtype = "html"
		email.Msg = "使用Golang发送邮件"
		email.Subject = "使用Golang发送邮件"
		fmt.Println(email)
		fmt.Println("Start Send msg", time.Now().Unix())

		emailerr := sendemail.SendMail(email) // send email
		if emailerr != nil {
			fmt.Println("send mail error!")
			fmt.Println(err)
		}else{
			fmt.Println("send mail success!")
		}

		erro := sendmsg.Sendmsg(msg, s) // send phone messages
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

func (p *DataPasser) MakeRender(w http.ResponseWriter, r *http.Request) {// 使用channel和并发

	r.ParseForm()
	getType := r.Header.Get("Content-Type")

	switch r.Method {
	case "POST":
		if getType != "application/x-www-form-urlencoded" {
			http.Error(w, "请使用application/x-www-form-urlencoded", 404)
			return
		}

		p.FormData <- r.Form
		var detail HttpRender
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

func (p *DataPasser) Log() { // 使用channel和并发
	for item := range p.FormData {
		msg := item.Get("msg")

		err := wechat.PostCustomMsg("farmer", msg) // Send wechat messages
		if err != nil {
			fmt.Println("Send wechat msg error")
		}
		MysqlErr := mysqldb.MysqlData(msg)	// save msg in db
		if MysqlErr != nil {
			fmt.Println("Save db Error")
		}
		s, e := mysqldb.MysqlSel()	// Get phone and email string from db
		var email sendemail.Email

		email.To = e
		email.Mailtype = "html"
		email.Msg = "使用Golang发送邮件"
		email.Subject = "使用Golang发送邮件"
		fmt.Println(email)
		fmt.Println("Start Send msg", time.Now().Unix())

		emailerr := sendemail.SendMail(email) // send email
		if emailerr != nil {
			fmt.Println("send mail error!")
			fmt.Println(err)
		}else{
			fmt.Println("send mail success!")
		}

		erro := sendmsg.Sendmsg(msg, s) // send phone messages
		if erro != nil {
			fmt.Println(erro.Error())
		}
		fmt.Println("End Send msg", time.Now().Unix())

		log.Println("发送成功！！")

		//fmt.Fprintf(w, string(body))

		fmt.Println("2. Item", item.Get("msg"))
	}
}
