package main

import (
	"log"
	"net/http"
	"./api"
	//"./email"

	"net/url"
	"runtime"
)


//主函数
func main() {
	log.Println("微信服务: 启动!")
	//email.SendEmail("明天出去游玩，各位带好东西")

	passer := &api.DataPasser{FormData: make(chan url.Values)}
	runtime.GOMAXPROCS(runtime.NumCPU())
	go passer.Log()

	http.HandleFunc("/api/", api.ApiRequest)
	http.HandleFunc("/wx/", api.Wxsend)
	http.HandleFunc("/send/", passer.MakeRender)
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		log.Fatal("微信服务: 端口监听和服务失败！, ", err)
	}
	//服务停止打印log
	log.Println("微信服务: 停止!")

}

