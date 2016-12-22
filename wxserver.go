package main

import (
	"log"
	"net/http"
	"./api"

)


//主函数
func main() {
	log.Println("微信服务: 启动!")
	http.HandleFunc("/api/", api.ApiRequest)
	http.HandleFunc("/wx/", api.Wxsend)
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		log.Fatal("微信服务: 端口监听和服务失败！, ", err)
	}
	//服务停止打印log
	log.Println("微信服务: 停止!")

}

