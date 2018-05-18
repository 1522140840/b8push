package main

import (
	"sync"
	"time"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
	"b8push/message"
	"github.com/fanliao/go-concurrentMap"
	"b8push/utils"
	"fmt"
)

var clientHandlers []func(*Client, message.Packet)

type Client struct {
	sync.Mutex
	Ws  *websocket.Conn
	ClientId    string
	topics *map[string]bool
	ConnectedAt int64
}

var udidClientMapping = concurrent.NewConcurrentMap(512, float32(0.7), 32)

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Ws: conn,
		ConnectedAt: time.Now().Unix(),
		ClientId:CreateClient(),

	}
}

func CreateClient() string{
	return uuid.NewV4().String()
}

func (c *Client) Start() {

	RegisterClient(c)

	log.Debugf("client++,clientid[%s]",c.ClientId)


}

func (client *Client)Read(){
	for {
		//设置超时时间  5min超时
		client.Ws.SetReadDeadline(time.Now().Add(5 * time.Minute))
		_, message, err := client.Ws.ReadMessage()
		if err != nil {
			log.Errorf("client read error,clientid:[%s]",client.ClientId)
			HandleDisconnect(client)
			break
		}
		log.Debugf("client read message,message:[%s]",message)

		if len(message)<4{
			log.Error("Read data length error")
			//断开连接
			HandleDisconnect(client)
		}
		client.handlePacket(message)
	}

}

func (c *Client) handlePacket(packetBytes []byte) {
	//保持和tcp push一致
	pckt, err := ParseClientPacket(packetBytes)
	if err != nil {
		log.Errorf("Packet error: %s. origin: %s", err.Error(), string(packetBytes))
		return
	}

	if len(clientHandlers)>pckt.Type()&&clientHandlers[pckt.Type()]!=nil{
		clientHandlers[pckt.Type()](c, pckt)
	}

}

func Write(client *Client,packet message.Packet){

	defer client.Unlock()

	client.Lock()
	data,err:=SerializeClientResult(packet)
	if err!=nil{
		log.Errorf("serializeClientResult error,error:[%s],clientId:[%s]",err,client.ClientId)
		return
	}
	err=client.Ws.WriteMessage(websocket.BinaryMessage, data)
	if err!=nil{
		log.Errorf("send data to client error,data:[%s],clientId:[%s]",data,client.ClientId)
	}else {
		log.Debugf("send data to client,data:[%s],clientId:[%s]",data,client.ClientId)
	}
}

/**
发送数据类型已经确认的二进制数据
 */
func WriteDateWithType(client *Client,data []byte,t int){

	defer client.Unlock()

	client.Lock()
	data,err:=SerializeClientDataAndTypeResult(data,t)
	if err!=nil{
		return
	}
	err=client.Ws.WriteMessage(websocket.BinaryMessage, data)
	if err!=nil{
		fmt.Println(err)
	}
}

func RegisterClient(c *Client) {
	 udidClientMapping.Put(c.ClientId, c)
}

func HandleDisconnect(c *Client){

	log.Debugf("client disconnect,clientId:[%s],topics:[%s]",c.ClientId,c.topics)
	udidClientMapping.Remove(c.ClientId)
	//删除订阅
	if c!=nil&&c.topics!=nil{
		for k,_:=range *c.topics{
			DelFromTopic(c,k)
		}
	}

	c.Ws.Close()
}

func init() {
	clientHandlers = make([]func(*Client, message.Packet), 200)
	clientHandlers[message.PACKAGE_SUBSCRIBE] = handleSubscribe
	clientHandlers[message.PACKET_PING] = handlePing

}

func handlePing(client *Client,p message.Packet){
	_, ok := p.(*message.Ping)
	if !ok{
		return
	}
	Write(client,PONG)
}

func handleSubscribe(client *Client,p message.Packet){
	result:=&message.SubscribeResult{
		Code:200,
	}

	subscribe, ok := p.(*message.Subscribe)
	if subscribe==nil||!ok{
		log.Errorf("Illegal  message")
		result.Code=201
		result.Msg="非法的订阅信息"
		return
	}
	//验证订阅内容是否合法
	infos,reslutStr:=CheckSubscribe(*subscribe)

	//订阅信息合法
	if util.StrIsBlank(reslutStr){
		resultTopics:= make(map[string]bool)

		for _,info:=range infos{
			topicName:=GetTopicNameBySubscribeInfo(info)
			if(!Add2Topic(client,info)){
				reslutStr+=topicName+";"
			}else{
				resultTopics[topicName]=true
			}

		}

		//删除无用订阅信息
		if client.topics!=nil {
			for k, _ := range *client.topics {
				if _, ok := resultTopics[k]; !ok {
					DelFromTopic(client, k)
				}
			}
		}
		client.topics=&resultTopics
	}else{
		log.Debugf("Illegal  message,[%s]",reslutStr)
	}

	if !util.StrIsBlank(reslutStr){
		result.Code=201
		result.Msg=reslutStr
		}

	Write(client,result)




}

