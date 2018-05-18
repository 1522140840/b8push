package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	//"time"
	//"time"
	"encoding/json"
	"flag"
)

var origin = "http://172.29.30.31:8080/"
var url = "ws://172.29.30.31:8080/echo"

//var origin = "http://127.0.0.1:8080/"
//var url = "ws://127.0.0.1:8080/echo"
type  Subscribe struct {
	Topics []SubscribeInfo "josn:`topics`"
}

type SubscribeInfo struct {
	Name string "json:`name`"			//名称
	Cycle string "json:`cycle`"		//周期
	DateType int "json:	`dataType`"   //订阅数据类型

}

var topic *string = flag.String("t", "huobi.btcusdt.kline", "topic")
var cycle *string = flag.String("c", "1min", "cycle")
var dataType *string = flag.String("d", "kline", "kline")
func main()  {
	flag.Parse()
	ws, err := websocket.Dial(url, "", "*")
	if err != nil {
		fmt.Println(err)
	}
	var data Subscribe
	if *dataType=="kline"{
		data=Subscribe{
			Topics:[]SubscribeInfo{SubscribeInfo{
				Name:*topic,
				Cycle:*cycle,
			}},
		}
	}else{
		data=Subscribe{
			Topics:[]SubscribeInfo{SubscribeInfo{
				Name:*topic,
			}},
		}
	}
	//取消订阅
	//data:=Subscribe{
	//	Topics:[]SubscribeInfo{},
	//}

	message,_ := json.Marshal(data)

	_, err = ws.Write(EnPackSendData(message,100))

	//ping := make([]byte, 4)
	//ping[0] = 0
	//ping[1] = 2
	//ping[2] = 1
	//ping[3] = 1
	//_, err = ws.Write(ping)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("Send: %s\n", message)

	for {
		var msg = make([]byte, 512)
		m, err := ws.Read(msg)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(msg[:m])
		fmt.Printf("Receive: %s\n", msg[:m])
	}


	ws.Close()//关闭连接


}

func EnPackSendData(sendBytes []byte, ty int) []byte {
	packetLength := len(sendBytes) + 4
	result := make([]byte, packetLength)
	result[0] = byte(((len(sendBytes) + 2) / 256))
	result[1] = byte(((len(sendBytes) + 2) % 256))
	result[2] = 1
	result[3] = byte(ty)
	copy(result[4:], sendBytes)

	//sendCrc := crc32.ChecksumIEEE(sendBytes)
	//result[packetLength-4] = byte(sendCrc >> 24)
	//result[packetLength-3] = byte(sendCrc >> 16 & 0xFF)
	//result[packetLength-2] = 0xFF
	//result[packetLength-1] = 0xFE
	//fmt.Println(result)
	//result[packetLength-1]=' '
	//fmt.Println(result)
	return result
}
