package main

import (
	"sync"
	"github.com/fanliao/go-concurrentMap"
	"b8push/message"
	"b8push/mysql"
	"github.com/robfig/cron"
	"b8push/utils"
	"time"
	"b8push/conf"
	"strings"

)


var topicsMap=concurrent.NewConcurrentMap(512, float32(0.7), 32)

/**
交易对资源  用于校验用户订阅的数据是否合法

 样例数据：symbols.Put("huobi.bchusdt.1min",1)   订阅信息-----数据类型【1：kline信息】
 */
var Symbols=concurrent.NewConcurrentMap(512,float32(0.7),32)

//从配置文件里取周期
var cicles []string


type TopicSlice struct {
	sync.Mutex
	TopicName string
	Type int
	messages chan message.Packet
	//在map和slice做的一个平衡，由于考虑到经常需要便利该集合进行消息分发，因此用slice
	Clients *[]*Client
}


func GetTopic(name string) (topic *TopicSlice){
	if v, ok := topicsMap.Get(name); ok==nil {
		if v!=nil{
			topic=v.(*TopicSlice)
		}

	}
	return topic
}


/**
添加订阅客户端
 */
func Add2Topic(client *Client,info message.SubscribeInfo) (flag bool ){

	name:=GetTopicNameBySubscribeInfo(info)

	if util.StrIsBlank(name){
		return
	}

	topic:=GetTopic(name)
	if topic==nil{
		log.Debugf("[%s]topic not exist",name)
		//创建topic
		topic,_=CreateTopic(info)
	}
	if topic!=nil {
		topic.Lock()
		x := append(*topic.Clients, client)
		topic.Clients = &x
		flag = true
		topic.Unlock()

		log.Debugf("[%s] add to [%s]topic success",client.ClientId,name)
		return
	}

	return
}
/**
删除订阅客户端
 */
 func DelFromTopic(client *Client,name string){
	 topic:=GetTopic(name)
	 if topic!=nil {

		 topic.Lock()
		 delIndex := -1
		 for i, c := range *topic.Clients {
			 if client.ClientId == c.ClientId {
				 delIndex = i
				 break
			 }
		 }
		 if delIndex>-1{
			 x:=append((*topic.Clients)[:delIndex],(*topic.Clients)[delIndex+1:]...)
			 topic.Clients=&x
		 }
		 topic.Unlock()
		 log.Debugf("remove [%s] from [%s]topic",client.ClientId,name)
	 }
 }


func ListenTopic(topic *TopicSlice){
	var (
		messageId int64
		sendRate  int=20//5条发一条
		limitFlag bool=false
		limitStartFlag bool=false
		sendCount int=0//一个周期已抛弃消息的数据
	)
	if topic.Type==message.RADAR_KLINE_DATA_TYPE{
		//只有kline数据限流
		limitFlag=true
	}

	log.Debugf("[%s] topic start listen",topic.TopicName)


	for{
		data:=<-topic.messages
		if limitFlag {
			//只对kline数据进行限流
			if len(topic.messages) > 100 {
				if !limitStartFlag {
					log.Debugf("[%s] topic limit rate start", topic.TopicName)
				}
				limitStartFlag = true

				kline := data.(*message.KLine)
				if messageId == kline.Tick.Id {
					//上一条消息id和当前消息id一致
					if sendCount == sendRate {
						sendCount = 0
					} else {
						//跳过本次推送
						sendCount += 1
						continue
					}

				} else {
					//上一条消息id和本条消息id不一致，直接发送数据
					//重新计数
					sendCount = 0
					messageId = kline.Tick.Id
				}
			}else{
				if limitStartFlag{
					log.Debugf("[%s] topic limit rate stop",topic.TopicName)
				}
				limitStartFlag=false
			}
		}

		for _,client:= range *topic.Clients{
			go Write(client,data)
		}
	}
}

func CreateTopic(info message.SubscribeInfo) (topic *TopicSlice,e error){

	name:=GetTopicNameBySubscribeInfo(info)

	if util.StrIsBlank(name){
		//非法的订阅消息
		log.Errorf("sub topic name can't null")
		return
	}

	if v, ok := topicsMap.Get(name); ok!=nil {
		log.Errorf("[%s] topic is not exist\n",name)
		topic=v.(*TopicSlice)
	} else {
		topic=&TopicSlice{
			TopicName:name,
			messages:make(chan message.Packet,200),
			Type:info.DateType,
			Clients:&[]*Client{},
		}
		_,e=topicsMap.Put(name,topic)
		if e!=nil{
			log.Errorf("[%s] topic create fail,error:[%s]",name,e)
		}else{
			//启动队列监听
			go ListenTopic(topic)

			//启动数据接收服务
			go ReciveDate(info)


			log.Debugf("[%s] topic create success",name)
		}

	}
	return  topic,e

}

func CheckSubscribe(subscribe message.Subscribe)([]message.SubscribeInfo,string){
	var (
		result string
		info =[]message.SubscribeInfo{}
	)
	for _,v:=range subscribe.Topics{
		k:=v.Name
		if !util.StrIsBlank(v.Cycle){
			k+="."+v.Cycle
		}
		t,_:=Symbols.Get(k)
		if t!=nil{
			info=append(info,message.SubscribeInfo{
				v.Name,
				v.Cycle,
				t.(int),
			})
		}else{
			//非法的订阅数据
			return nil,k+":非法的订阅信息"

		}
	}

	return info,result
}


func init() {
	log.Debugf("init load Symbols and Exchanges")
	//加载周期
	cicleStr:=conf.GetVal("symbol","cycle")
	if !util.StrIsBlank(cicleStr){
		cicles=strings.Split(cicleStr,",")
	}


	//加载交易对
	InitSymbol()
	symbol := cron.New()
	symbolSpec := "0 0 0 * * ?"
	symbol.AddFunc(symbolSpec,InitSymbol)
	symbol.Start()

	//加载交易所
	InitExchange()
	exchange := cron.New()
	exchangeSpec := "0 0 0 * * ?"
	exchange.AddFunc(exchangeSpec,InitExchange)
	exchange.Start()


}

func GetTopicNameBySubscribeInfo(info message.SubscribeInfo)string{
	var result string
	if info!=(message.SubscribeInfo{})&&!util.StrIsBlank(info.Name){
		if util.StrIsBlank(info.Cycle){
			result=info.Name
		}else{
			result=info.Name+"."+info.Cycle
		}

	}

	return result
}

//初始化所有交易对
func InitSymbol(){

	for{
		symbols,error:=mysql.QuerySymbol()
		if error==nil&&len(symbols)>0{
			for k,_:=range symbols{
				for _,c:=range  cicles{
					Symbols.Put(k+".kline."+c,message.RADAR_KLINE_DATA_TYPE)
				}
			}
			log.Debugf("初始化【%d】个交易对\n",len(symbols))
			break
		}else{
			//重连
			log.Errorf("获取交易对异常，3S 后进行重连。error:",error)
			time.Sleep(time.Second*3)
		}

	}

}

//初始化所有交易所
func InitExchange(){

	for{
		exchanges,error:=mysql.QueryExchange()
		if error==nil&&len(exchanges)>0{
			for k,_:=range exchanges{
				Symbols.Put(k,message.RADAR_24_HOUR_DATA_TYPE)
			}

			log.Debugf("初始化【%d】个交易所\n",len(exchanges))
			break
		}else{
			//重连
			log.Errorf("获取交易所异常，3S 后进行重连。error:",error)
			time.Sleep(time.Second*3)
		}

	}

}
