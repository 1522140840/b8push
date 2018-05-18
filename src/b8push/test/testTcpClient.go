package main


import (

"encoding/json"

	"net"
)

type  Subscribe struct {
	Topics string "josn:`topics`"
}

func main()  {
	hawkServer, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

	//连接服务器
	connection, _ := net.DialTCP("tcp", nil, hawkServer)

	data:=Subscribe{
		Topics:"huobi_1min",
	}

	message,_ := json.Marshal(data)
	//for{
	connection.Write(EnPackSendData(message,200))

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



