package message

//数据类型从100开始避免与Tcp push消息类型冲突
const(
	PACKAGE_SUBSCRIBE =100+iota
	PACKAGE_KLINE
	PACKAGE_EXCHANGE24HOUR

)

const (
	RADAR_KLINE_DATA_TYPE=1+iota   //交易对kline数据
	RADAR_24_HOUR_DATA_TYPE
)



const (
	PACKET_PING=1
	PACKET_PONG=2
)

var PONG =&Pong{}

type Packet interface {
	TypeString() string
	Type() int
}



//KLine信息
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

func (exchange24Hour *Exchange24Hour)TypeString() string {
	return "exchange24Hour"
}

func (exchange24Hour *Exchange24Hour)Type() int {
	return PACKAGE_EXCHANGE24HOUR
}

func (kline *KLine)TypeString() string {
	return "kline"
}

func (kline *KLine)Type() int {
	return PACKAGE_KLINE
}

/**
订阅消息
 */
type  Subscribe struct {
	Topics []SubscribeInfo "josn:`topics`"
}
type SubscribeInfo struct {
	Name string "json:`name`"			//名称
	Cycle string "json:`cycle`"		//周期
	DateType int "json:	`dataType`"   //订阅数据类型

}


func (sub *Subscribe)TypeString() string {
	return "subscribe"
}

func (sub *Subscribe)Type() int {
	return PACKAGE_SUBSCRIBE
}

type SubscribeResult struct{
	Code          int    `json:"code"`
	Msg           string `json:"msg"`
}

func (p *SubscribeResult) TypeString() string {
	return "subscriberesult"
}
func (p *SubscribeResult) Type() int {
	return PACKAGE_SUBSCRIBE
}

/////////////////////////////// Ping //////////////////////////////////////////////////

type Ping struct{}

func (p *Ping) TypeString() string { return "ping" }

func (p *Ping) Type() int { return PACKET_PING }

func (p *Ping) SerializeForClient() ([]byte, error) {
	return nil, nil //FIXME 使用常量
}

func (p *Ping) ReceivedAt() int64 { return 0 }

/////////////////////////////// Pong //////////////////////////////////////////////////

type Pong struct{}

func (p *Pong) TypeString() string { return "pong" }

func (p *Pong) Type() int { return PACKET_PONG }

func (p *Pong) SerializeForClient() ([]byte, error) {
	return nil, nil
}

func (p *Pong) ReceivedAt() int64 { return 0 }

