package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"github.com/gin-gonic/gin/json"
	"math/rand"
	"strconv"
	"strings"
)
type KLine struct {
	Ch string 	`json:"ch"`  	//订阅串 "huobi.btcusdt.kline.1min",
	Ts int64 	`json:"ts"`	//时间戳
	Tick Tick 	`json:"tick"`	//KLine详细信息
}
//KLine详细信息
type Tick struct {
	Id int64 		`json:"id"`  		//消息ID
	Open float32 	`json:"open"` 		//开盘价格
	Close float32 	`json:"close"` 	//收盘价格
	Low float32 	`json:"low"` 		//最低价格
	High float32 	`json:"high"` 		//最高价格
	Vol float32 	`json:"vol"`		//成交额, 即 sum(每一笔成交价 * 该笔的成交量)
}

type Exchange24Hour struct {
	Ch string 	`json:"ch"`  	//huobi.overview
	Ts int64 	`json:"ts"`	//时间戳
	Status string `json:"status"`	//状态
	Data []ExchangeDetail `json:"data"` //详情
}

type ExchangeDetail struct {
	Open float32 	`json:"open"` 		//开盘价格
	Close float32 	`json:"close"` 	//收盘价格
	Low float32 	`json:"low"` 		//最低价格
	High float32 	`json:"high"` 		//最高价格
	Amount float32 `json:"amount"`
	Vol float32 	`json:"vol"`		//成交额, 即 sum(每一笔成交价 * 该笔的成交量)
	Count int32 	`json:"count"`
	Symbol string 	`json:"symbol"`	//交易对
}

type sub_result struct {
	Status string `json:"status"`
	Subbed string `json:"subbed"`
	Ts  int `json:"ts"`
}

var upgrader = websocket.Upgrader{}
func main() {
	//http.HandleFunc("/echo", websocket.Handler())
	http.HandleFunc("/", radarHandler)

	err := http.ListenAndServe("127.0.0.1:9002", nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func  radarHandler(w http.ResponseWriter, r *http.Request)  {
	c, _ := upgrader.Upgrade(w, r, nil)

	fmt.Println("clent++")

	go read(c)


}

func read(conn *websocket.Conn){
	for {

		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("read error:%s",err)
			break
		}
		fmt.Printf("message:%s\n",message)
		data,_:=json.Marshal(sub_result{
			Status:"ok",
			Subbed:string(message),
			Ts:112312,
		})

		fmt.Println(data)

		conn.WriteMessage(websocket.BinaryMessage,data)

		if strings.Contains(string(message),"overview"){
			go write24Hour(conn,string(message))
		}else {
			go writeKline(conn,string(message))
		}


	}
}

func writeKline(conn *websocket.Conn,topic string){


		for{
			msg:=KLine{
				Ch:topic+".kline."+strconv.Itoa(rand.Intn(2))+"min",
				Ts:112312312,
				Tick:Tick{
					Id:1231,
					Open:12.1,
					Close:12.3,
					Low:18.1,
					High:902.1,
					Vol:78.1,
				},
			}
			data,_:=json.Marshal(msg)

			conn.WriteMessage(websocket.BinaryMessage,data)

			//time.Sleep(time.Second)
		}
}

func write24Hour(conn *websocket.Conn,topic string){


	for{
		msg:=Exchange24Hour{
			Ch:topic,
			Ts:112312312,
			Status:"ok",
			Data:[]ExchangeDetail{
				ExchangeDetail{
					Open :123.123,
					Close:12.1,
					Low:900.1,
					High:19.1,
					Amount:21.1,
					Vol:19023.1,
					Count:1091,
					Symbol:"eosusdt",
				},
				ExchangeDetail{
					Open :123.123,
					Close:12.1,
					Low:900.1,
					High:19.1,
					Amount:21.1,
					Vol:19023.1,
					Count:1091,
					Symbol:"ethusdt",
				},
			},
		}
		data,_:=json.Marshal(msg)
		conn.WriteMessage(websocket.BinaryMessage,data)

		//time.Sleep(time.Second)
	}
}
