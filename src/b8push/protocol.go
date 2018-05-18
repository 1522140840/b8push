package main

import (
	"encoding/base64"
	"b8push/message"
	"encoding/json"
	"bytes"
	"encoding/binary"
)

const (
	PACKET_VER_1 = byte(1)
)

var (
	PONG = &message.Pong{}
	PING = &message.Ping{}
)

func ParseClientPacket(data []byte) (message.Packet, error) {

	if len(data)<4{
		//一条消息至少四个字节
		return nil, newProtocolErr("packet info error")
	}

	/**
	兼容tcp push
	 */
	pVer := data[2]
	pType := data[3]
	if len(data)>4{
		data = data[4:]
	}


	if pVer != PACKET_VER_1 {
		return nil, newProtocolErr("Unsupported protocol ver")
	}

	if pType == 21 {
		log.Debugf("client package version:%d,type:%d,size:%d,data:%s", pVer, pType, len(data), base64.StdEncoding.EncodeToString(data))
	} else {
		log.Debugf("client package version:%d,type:%d,size:%d,data:%s", pVer, pType, len(data), string(data))
	}

	if pType == message.PACKET_PING {
		return PING, nil
	}

	switch pType {
	case message.PACKAGE_SUBSCRIBE:
		var p *message.Subscribe
		err := json.Unmarshal(data, &p)
		if err != nil {
			return nil, err
		}
		return p, nil

	default:
		return nil, newProtocolErr("Unknown packet type")
	}

}

func SerializeClientResult(packet message.Packet) ([]byte,error){

	var result bytes.Buffer

	var b []byte
	switch packet.Type() {
	case message.PACKAGE_SUBSCRIBE:
		//订阅结果
		subscribeResult:=packet.(*message.SubscribeResult)
		b,_=json.Marshal(subscribeResult)

		err:= binary.Write(&result, binary.BigEndian, int16(len(b) + 2))
		if err != nil {
			return nil, err
		}
	case message.PACKAGE_KLINE:
		//kline数据
		kline:=packet.(*message.KLine)

		b,_=json.Marshal(kline)

		err:= binary.Write(&result, binary.BigEndian, int16(len(b) + 2))
		if err != nil {
			return nil, err
		}
	case message.PACKET_PONG:
		err:= binary.Write(&result, binary.BigEndian, int16(2))
		if err != nil {
			return nil, err
		}
	case message.PACKAGE_EXCHANGE24HOUR:
		//交易所24小时数据
		exchange24Hour:=packet.(*message.Exchange24Hour)

		b,_=json.Marshal(exchange24Hour)

		err:= binary.Write(&result, binary.BigEndian, int16(len(b) + 2))
		if err != nil {
			return nil, err
		}
	}

	binary.Write(&result, binary.BigEndian, PACKET_VER_1)

	err:= binary.Write(&result, binary.BigEndian, int8(packet.Type())) // code=20 result
	if err != nil {
		return nil, err
	}
	result.Write(b)

	return result.Bytes(), nil

}

func SerializeClientDataAndTypeResult(data []byte,t int) ([]byte,error){

	var result bytes.Buffer


	err:= binary.Write(&result, binary.BigEndian, int16(len(data) + 2))
	if err != nil {
		return nil, err
	}

	binary.Write(&result, binary.BigEndian, PACKET_VER_1)

	err= binary.Write(&result, binary.BigEndian, int8(t)) // code=20 result
	if err != nil {
		return nil, err
	}
	result.Write(data)

	return result.Bytes(), nil

}



type ProtocolError struct {
	Desc string
}

func (e *ProtocolError) Error() string {
	return e.Desc
}

func newProtocolErr(desc string) *ProtocolError {
	return &ProtocolError{Desc: desc}
}
