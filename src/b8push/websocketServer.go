package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"b8push/conf"
	"b8push/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//允许跨域调用
		return true
	},
}

func StartServer(){

	port:=conf.GetVal("b8push","websocket_server_port")

	//默认8080端口
	if util.StrIsBlank(port){
		port="8080"
	}

	http.HandleFunc("/echo", echoHandler)

	err := http.ListenAndServe(":"+port, nil)

	log.Infof("start server,server port:[%s]",port)

	if err != nil {
		log.Errorf("server start error:%s",err)
		panic("ListenAndServe: " + err.Error())
	}

}

func  echoHandler(w http.ResponseWriter, r *http.Request)  {
	c, err := upgrader.Upgrade(w, r, nil)
	if err!=nil{
		log.Errorf("client connect server error:[%s]",err)
		return
	}
	client:=NewClient(c)
	client.Start()

	go client.Read()


}


