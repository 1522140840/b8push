package main

import (
	"fmt"
	"time"
	"github.com/gorilla/websocket"
	"flag"
	ul "net/url"
)

type sub_kline struct {
	Sub string `json:sub`
}

type sub_result struct {
	status string
	subbed string
	ts  int
}

/**
获取雷达数据
 */

var addr = flag.String("addr", "10.103.197.205:9002", "http service address")

func main(){
	topic:=[]string{"huobi.bchusdt","huobi.ethusdt","huobi.etcusdt","huobi.ltcusdt","huobi.eosusdt","huobi.xrpusdt","huobi.omgusdt","huobi.dashusdt","huobi.zecusdt","huobi.adausdt"}
	for _,v:=range topic{
		for i:=0;i<1;i++ {
			go ReciveDate(v, 1)
		}
	}


	time.Sleep(time.Minute*60)
}

func ReciveDate(sub string,t int){

	u := ul.URL{Scheme: "ws", Host: *addr, Path: "/"}
	var dialer *websocket.Dialer

	con, _, err := dialer.Dial(u.String(), nil)

	if err!=nil{
		//关闭连接通知 topic进行重连操作
	}

	con.WriteMessage(websocket.BinaryMessage,[]byte(sub))

	for {
		_, mess, err :=con.ReadMessage()
		if err!=nil{
			fmt.Printf("read radar data error,err:%s \n", err)
			//重连
			//。。。
			break
		}
		fmt.Printf("message: %s\n",mess)
	}



}
