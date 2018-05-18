package main

import (
	"b8push/message"
	al "b8push/log"
	"time"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
	"github.com/fanliao/go-concurrentMap"
	"strings"
	"b8push/conf"
	"github.com/cihub/seelog"
)


type sub_result struct {
	Status string `json:"status"`
	Subbed string `json:"subbed"`
	Ts  int `json:"ts"`
}

/**
获取雷达数据
 */
var addr string

var RadarChannel=concurrent.NewConcurrentMap(512,float32(0.7),32)

var radarLog seelog.LoggerInterface

/**
启动数据接收服务
 */
func ReciveDate(info message.SubscribeInfo){

	v,_:=RadarChannel.Get(info.Name)
	if v!=nil{
		//不需要重复创建
		return
	}

	RadarChannel.Put(info.Name,1)

	 u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	 var dialer *websocket.Dialer

	 con, _, err := dialer.Dial(u.String(), nil)

	 if err!=nil{
	 	radarLog.Errorf("连接radar server异常,error:%s\n",err)

		go ReconectRadar(info)
		return
	 }

	 if subData(con,info){
		 go hanlerClient(con,info)

	 }else{
	 	//订阅失败重连
		go ReconectRadar(info)
	 }


 }
 /**
 订阅数据
  */
 func subData(con *websocket.Conn,info message.SubscribeInfo)(flag bool){

 	sub:=info.Name
 	if info.DateType==1{
		sub=strings.Split(sub,".kline")[0]
	}

	 con.WriteMessage(websocket.BinaryMessage,[]byte(sub))

	 //接收订阅返回结果
	 _, mess, err :=con.ReadMessage()

	 if err!=nil{
		 radarLog.Errorf("订阅消息接收失败,[%s]",err)
		 return flag
	 }

	 radarLog.Debugf("订阅结果:%s\n",mess)

	 result:=&sub_result{}
	 err=json.Unmarshal(mess,result)

	 if err!=nil{
		 radarLog.Errorf("订阅消息接收解析失败,[%s]",err)
		 return flag
	 }else{
		if result.Status=="ok"{
			return true
		}else {
			return flag
		}
	 }



 }

 func hanlerClient( con *websocket.Conn,info message.SubscribeInfo){

 	for {
 		//设置超时时间  10S超时
		con.SetReadDeadline(time.Now().Add(10 * time.Second))
		_, mess, err :=con.ReadMessage()
		if err!=nil{
			radarLog.Debugf("read radar data error,err:%s \n", err)
			//重连
			go ReconectRadar(info)
			break
		}
		switch info.DateType {
			case message.RADAR_KLINE_DATA_TYPE://kline数据
				pack:=&message.KLine{}
				err=json.Unmarshal(mess,pack)
				if err==nil{
					t:=GetTopic(pack.Ch)
					if t!=nil{
						t.messages<-pack
						radarLog.Debugf("read message from radar,topic[%s],message[%s]",t.TopicName,pack)
					}
				}
			case message.RADAR_24_HOUR_DATA_TYPE:
				exhange:=&message.Exchange24Hour{}
				err=json.Unmarshal(mess,exhange)
				if err==nil{
					t:=GetTopic(exhange.Ch)
					if t!=nil{
						t.messages<-exhange
						radarLog.Debugf("read message from radar,topic[%s],message[%s]",t.TopicName,exhange)
					}
				}
		}



	}
 }

 func ReconectRadar(info message.SubscribeInfo){

	 RadarChannel.Remove(info.Name)

 	//避免不间断的重连
 	time.Sleep(time.Second*3)

 	radarLog.Debugf("radar开始重连，message:%s\n",info)

 	go ReciveDate(info)

 }

func  TestReciveDate(info message.SubscribeInfo)  {
	v,_:=RadarChannel.Get(info.Name)
	if v!=nil{
		//不需要重复创建
		return
	}

	RadarChannel.Put(info.Name,1)

	msg:=message.KLine{
		Ch:"huobi.btcusdt.kline.1min",
		Ts:112312312,
		Tick:message.Tick{
			Id:1231,
			Open:12.1,
			Close:12.3,
			Low:18.1,
			High:902.1,
			Vol:78.1,
		},
	}
	for{

		t:=GetTopic(msg.Ch)
		if t!=nil{
			t.messages<-&msg
		}
		time.Sleep(time.Second)
	}

}

func init(){
	radarLog=al.InitLog("radar")

	server:=conf.GetVal("radar","server")

	port:=conf.GetVal("radar","port")

	//不能过分强依赖于雷达数据
	//if util.StrIsBlank(server)||util.StrIsBlank(port){
	//	panic("radar dataSource not config")
	//}
	addr=server+":"+port

	radarLog.Debugf("init radar info，server:[%s]",addr)


}




